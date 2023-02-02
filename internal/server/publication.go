package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) GetPublication(ctx context.Context, req *api.GetPublicationRequest) (*api.GetPublicationResponse, error) {
	pub, err := s.services.Repository.GetPublication(req.Id)
	if err != nil {
		// TODO How do we differentiate between errors? e.g. NotFound vs. Internal (database unavailable,...)
		return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %w", req.Id, err)
	}

	res := &api.GetPublicationResponse{Publication: PublicationToMessage(pub)}

	return res, nil
}

func (s *server) GetAllPublications(req *api.GetAllPublicationsRequest, stream api.Biblio_GetAllPublicationsServer) (err error) {
	// TODO errors in EachPublication aren't caught and pushed upstream. Returning 'false' in the callback
	//   breaks the loop, but EachPublication will return 'nil'.
	//
	//   Logging during streaming doesn't work / isn't possible. The grpc_zap interceptor is only called when
	// 	 GetAllPublication returns an error.
	errorStream := s.services.Repository.EachPublication(func(p *models.Publication) bool {
		res := &api.GetAllPublicationsResponse{Publication: PublicationToMessage(p)}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get all publications: %s", errorStream)
	}

	return nil
}

func (s *server) SearchPublications(ctx context.Context, req *api.SearchPublicationsRequest) (*api.SearchPublicationsResponse, error) {
	page := 1
	if req.Limit > 0 {
		page = int(req.Offset)/int(req.Limit) + 1
	}
	args := models.NewSearchArgs().WithQuery(req.Query).WithPage(page)
	hits, err := s.services.PublicationSearchService.Search(args)
	if err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "Could not search publications: %s :: %s", req.Query, err)
	}

	res := &api.SearchPublicationsResponse{
		Limit:  int32(hits.Limit),
		Offset: int32(hits.Offset),
		Total:  int32(hits.Total),
		Hits:   make([]*api.Publication, len(hits.Hits)),
	}
	for i, p := range hits.Hits {
		res.Hits[i] = PublicationToMessage(p)
	}

	return res, nil
}

func (s *server) UpdatePublication(ctx context.Context, req *api.UpdatePublicationRequest) (*api.UpdatePublicationResponse, error) {
	p := MessageToPublication(req.Publication)

	// TODO Fetch user information via better authentication (no basic auth)
	user := &models.User{
		ID:       "n/a",
		FullName: "system user",
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for publication %s: %s", p.ID, err)
	}

	if err := s.services.Repository.UpdatePublication(req.Publication.SnapshotId, p, user); err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "failed to store publication %s, %s", p.ID, err)
	}
	if err := s.services.PublicationSearchService.Index(p); err != nil {
		// TODO How do we differentiate between errors
		return nil, status.Errorf(codes.Internal, "failed to index publication %s, %s", p.ID, err)
	}

	return &api.UpdatePublicationResponse{}, nil
}

// TODO catch indexing errors
func (s *server) AddPublications(stream api.Biblio_AddPublicationsServer) error {
	// indexing channel
	indexC := make(chan *models.Publication)

	var indexWG sync.WaitGroup

	// start bulk indexer
	indexWG.Add(1)
	go func() {
		defer indexWG.Done()
		s.services.PublicationSearchService.IndexMultiple(indexC)
	}()

	defer func() {
		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	}()

	var seq int

	for {
		seq++

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read stream: %s", err)
		}

		p := MessageToPublication(req.Publication)

		if p.ID == "" {
			p.ID = ulid.Make().String()
		}
		if p.Status == "" {
			p.Status = "private"
		}
		if p.Classification == "" {
			p.Classification = "U"
		}
		for i, val := range p.Abstract {
			if val.ID == "" {
				val.ID = ulid.Make().String()
			}
			p.Abstract[i] = val
		}
		for i, val := range p.LaySummary {
			if val.ID == "" {
				val.ID = ulid.Make().String()
			}
			p.LaySummary[i] = val
		}
		for i, val := range p.Link {
			if val.ID == "" {
				val.ID = ulid.Make().String()
			}
			p.Link[i] = val
		}

		// TODO this should return structured messages (see validate)
		if err := p.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
		if err := s.services.Repository.SavePublication(p, nil); err != nil {
			msg := fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
				return status.Errorf(codes.Internal, msg)
			}
			continue
		}

		indexC <- p
	}
}

// TODO catch indexing errors
func (s *server) ImportPublications(stream api.Biblio_ImportPublicationsServer) error {
	// indexing channel
	indexC := make(chan *models.Publication)

	var indexWG sync.WaitGroup

	// start bulk indexer
	indexWG.Add(1)
	go func() {
		defer indexWG.Done()
		s.services.PublicationSearchService.IndexMultiple(indexC)
	}()

	defer func() {
		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	}()

	var seq int

	for {
		seq++

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read stream: %s", err)
		}

		p := MessageToPublication(req.Publication)

		// TODO this should return structured messages (see validate)
		if err := p.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.ImportPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
		if err := s.services.Repository.ImportCurrentPublication(p); err != nil {
			msg := fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.ImportPublicationsResponse{Message: msg}); err != nil {
				return status.Errorf(codes.Internal, msg)
			}
			continue
		}

		indexC <- p
	}
}

func (s *server) GetPublicationHistory(req *api.GetPublicationHistoryRequest, stream api.Biblio_GetPublicationHistoryServer) (err error) {
	errorStream := s.services.Repository.PublicationHistory(req.Id, func(p *models.Publication) bool {
		res := &api.GetPublicationHistoryResponse{Publication: PublicationToMessage(p)}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get publication history: %s", errorStream)
	}

	return nil
}

func (s *server) PurgePublication(ctx context.Context, req *api.PurgePublicationRequest) (*api.PurgePublicationResponse, error) {
	if err := s.services.Repository.PurgePublication(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication with id %s: %s", req.Id, err)
	}
	if err := s.services.PublicationSearchService.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication from index with id %s: %s", req.Id, err)
	}

	return &api.PurgePublicationResponse{}, nil
}

func (s *server) PurgeAllPublications(ctx context.Context, req *api.PurgeAllPublicationsRequest) (*api.PurgeAllPublicationsResponse, error) {
	if err := s.services.Repository.PurgeAllPublications(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all publications: %s", err)
	}
	// TODO use delete by query instead of recreating?
	if err := s.services.PublicationSearchService.DeleteIndex(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete publication index: %s", err)
	}
	if err := s.services.PublicationSearchService.CreateIndex(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not create publication index: %s", err)
	}

	return &api.PurgeAllPublicationsResponse{}, nil
}

func (s *server) ValidatePublications(stream api.Biblio_ValidatePublicationsServer) error {
	var seq int32

	for {
		seq++

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read stream: %s", err)
		}

		p := MessageToPublication(req.Publication)

		err = p.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			res := &api.ValidatePublicationsResponse{Seq: seq, Id: p.ID, Message: validationErrs.Error()}
			if err = stream.Send(res); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}

func PublicationToMessage(p *models.Publication) *api.Publication {
	msg := &api.Publication{}

	msg.Id = p.ID

	switch p.PublicationStatus {
	case "unpublished":
		msg.PublicationStatus = api.Publication_PUBLICATION_STATUS_UNPUBLISHED
	case "accepted":
		msg.PublicationStatus = api.Publication_PUBLICATION_STATUS_ACCEPTED
	case "published":
		msg.PublicationStatus = api.Publication_PUBLICATION_STATUS_PUBLISHED
	}

	switch p.Type {
	case "journal_article":
		msg.Type = api.Publication_TYPE_JOURNAL_ARTICLE
	case "book":
		msg.Type = api.Publication_TYPE_BOOK
	case "book_chapter":
		msg.Type = api.Publication_TYPE_BOOK_CHAPTER
	case "book_editor":
		msg.Type = api.Publication_TYPE_BOOK_EDITOR
	case "issue_editor":
		msg.Type = api.Publication_TYPE_ISSUE_EDITOR
	case "conference":
		msg.Type = api.Publication_TYPE_CONFERENCE
	case "dissertation":
		msg.Type = api.Publication_TYPE_DISSERTATION
	case "miscellaneous":
		msg.Type = api.Publication_TYPE_MISCELLANEOUS
	}

	switch p.Status {
	case "private":
		msg.Status = api.Publication_STATUS_PRIVATE
	case "public":
		msg.Status = api.Publication_STATUS_PUBLIC
	case "deleted":
		msg.Status = api.Publication_STATUS_DELETED
	case "returned":
		msg.Status = api.Publication_STATUS_RETURNED
	}

	msg.HasBeenPublic = p.HasBeenPublic

	for _, val := range p.Abstract {
		msg.Abstract = append(msg.Abstract, &api.Text{
			Id:   val.ID,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	msg.AdditionalInfo = p.AdditionalInfo

	msg.AlternativeTitle = p.AlternativeTitle

	msg.ArticleNumber = p.ArticleNumber

	msg.ArxivId = p.ArxivID

	for _, val := range p.Author {
		var depts []*api.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, &api.ContributorDepartment{
				Id:   dept.ID,
				Name: dept.Name,
			})
		}

		msg.Author = append(msg.Author, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	msg.BatchId = p.BatchID

	switch p.Classification {
	case "U":
		msg.Classification = api.Publication_CLASSIFICATION_U
	case "A1":
		msg.Classification = api.Publication_CLASSIFICATION_A1
	case "A2":
		msg.Classification = api.Publication_CLASSIFICATION_A2
	case "A3":
		msg.Classification = api.Publication_CLASSIFICATION_A3
	case "A4":
		msg.Classification = api.Publication_CLASSIFICATION_A4
	case "B1":
		msg.Classification = api.Publication_CLASSIFICATION_B1
	case "B2":
		msg.Classification = api.Publication_CLASSIFICATION_B2
	case "B3":
		msg.Classification = api.Publication_CLASSIFICATION_B3
	case "C1":
		msg.Classification = api.Publication_CLASSIFICATION_C1
	case "C3":
		msg.Classification = api.Publication_CLASSIFICATION_C3
	case "D1":
		msg.Classification = api.Publication_CLASSIFICATION_D1
	case "D2":
		msg.Classification = api.Publication_CLASSIFICATION_D2
	case "P1":
		msg.Classification = api.Publication_CLASSIFICATION_P1
	case "V":
		msg.Classification = api.Publication_CLASSIFICATION_V
	}

	if p.DateCreated != nil {
		msg.DateCreated = timestamppb.New(*p.DateCreated)
	}
	if p.DateUpdated != nil {
		msg.DateUpdated = timestamppb.New(*p.DateUpdated)
	}
	if p.DateFrom != nil {
		msg.DateFrom = timestamppb.New(*p.DateFrom)
	}
	if p.DateUntil != nil {
		msg.DateUntil = timestamppb.New(*p.DateUntil)
	}

	msg.Extern = p.Extern

	msg.Title = p.Title

	msg.DefensePlace = p.DefensePlace
	msg.DefenseDate = p.DefenseDate
	msg.DefenseTime = p.DefenseTime

	msg.ConferenceName = p.ConferenceName
	msg.ConferenceLocation = p.ConferenceLocation
	msg.ConferenceOrganizer = p.ConferenceOrganizer
	msg.ConferenceStartDate = p.ConferenceStartDate
	msg.ConferenceEndDate = p.ConferenceEndDate

	for _, val := range p.Department {
		var depts []*api.DepartmentRef

		for _, dept := range val.Tree {
			depts = append(depts, &api.DepartmentRef{
				Id: dept.ID,
			})
		}

		msg.Department = append(msg.Department, &api.Department{
			Id:   val.ID,
			Tree: depts,
		})
	}

	if p.Creator != nil {
		msg.Creator = &api.RelatedUser{Id: p.Creator.ID, Name: p.Creator.Name}
	}

	if p.User != nil {
		msg.User = &api.RelatedUser{Id: p.User.ID, Name: p.User.Name}
	}

	if p.LastUser != nil {
		msg.LastUser = &api.RelatedUser{Id: p.LastUser.ID, Name: p.LastUser.Name}
	}

	msg.Doi = p.DOI

	msg.Edition = p.Edition

	for _, val := range p.Editor {
		var depts []*api.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, &api.ContributorDepartment{
				Id:   dept.ID,
				Name: dept.Name,
			})
		}
		msg.Editor = append(msg.Editor, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
			CreditRole: val.CreditRole,
		})
	}

	msg.Eisbn = p.EISBN

	msg.Eissn = p.EISSN

	msg.EsciId = p.ESCIID

	for _, val := range p.File {
		f := &api.File{
			Id:           val.ID,
			License:      val.License,
			ContentType:  val.ContentType,
			EmbargoDate:  val.EmbargoDate,
			Name:         val.Name,
			Size:         int32(val.Size),
			Sha256:       val.SHA256,
			OtherLicense: val.OtherLicense,
		}

		switch val.AccessLevel {
		case "info:eu-repo/semantics/openAccess":
			f.AccessLevel = api.File_ACCESS_LEVEL_OPEN_ACCESS
		case "info:eu-repo/semantics/embargoedAccess":
			f.AccessLevel = api.File_ACCESS_LEVEL_EMBARGOED_ACCESS
		case "info:eu-repo/semantics/restrictedAccess":
			f.AccessLevel = api.File_ACCESS_LEVEL_RESTRICTED_ACCESS
		case "info:eu-repo/semantics/closedAccess":
			f.AccessLevel = api.File_ACCESS_LEVEL_CLOSED_ACCESS
		}

		if val.DateCreated != nil {
			f.DateCreated = timestamppb.New(*val.DateCreated)
		}
		if val.DateUpdated != nil {
			f.DateUpdated = timestamppb.New(*val.DateUpdated)
		}

		switch val.AccessLevelDuringEmbargo {
		case "info:eu-repo/semantics/closedAccess":
			f.AccessLevelDuringEmbargo = api.File_ACCESS_LEVEL_CLOSED_ACCESS
		case "info:eu-repo/semantics/restrictedAccess":
			f.AccessLevelDuringEmbargo = api.File_ACCESS_LEVEL_RESTRICTED_ACCESS
		}

		switch val.AccessLevelAfterEmbargo {
		case "info:eu-repo/semantics/openAccess":
			f.AccessLevelAfterEmbargo = api.File_ACCESS_LEVEL_OPEN_ACCESS
		case "info:eu-repo/semantics/restrictedAccess":
			f.AccessLevelAfterEmbargo = api.File_ACCESS_LEVEL_RESTRICTED_ACCESS
		}

		switch val.PublicationVersion {
		case "publishedVersion":
			f.PublicationVersion = api.File_PUBLICATION_VERSION_PUBLISHED_VERSION
		case "authorVersion":
			f.PublicationVersion = api.File_PUBLICATION_VERSION_AUTHOR_VERSION
		case "acceptedVersion":
			f.PublicationVersion = api.File_PUBLICATION_VERSION_ACCEPTED_VERSION
		case "updatedVersion":
			f.PublicationVersion = api.File_PUBLICATION_VERSION_UPDATED_VERSION
		}

		switch val.Relation {
		case "main_file":
			f.Relation = api.File_RELATION_MAIN_FILE
		case "colophon":
			f.Relation = api.File_RELATION_COLOPHON
		case "data_fact_sheet":
			f.Relation = api.File_RELATION_DATA_FACT_SHEET
		case "peer_review_report":
			f.Relation = api.File_RELATION_PEER_REVIEW_REPORT
		case "table_of_contents":
			f.Relation = api.File_RELATION_TABLE_OF_CONTENTS
		case "agreement":
			f.Relation = api.File_RELATION_AGREEMENT
		case "supplementary_material":
			f.Relation = api.File_RELATION_SUPPLEMENTARY_MATERIAL
		}

		msg.File = append(msg.File, f)
	}

	msg.Handle = p.Handle

	switch p.HasConfidentialData {
	case "yes":
		msg.HasConfidentialData = api.Confirmation_CONFIRMATION_YES
	case "no":
		msg.HasConfidentialData = api.Confirmation_CONFIRMATION_NO
	case "dontknow":
		msg.HasConfidentialData = api.Confirmation_CONFIRMATION_DONT_KNOW
	}
	switch p.HasPatentApplication {
	case "yes":
		msg.HasPatentApplication = api.Confirmation_CONFIRMATION_YES
	case "no":
		msg.HasPatentApplication = api.Confirmation_CONFIRMATION_NO
	case "dontknow":
		msg.HasPatentApplication = api.Confirmation_CONFIRMATION_DONT_KNOW
	}
	switch p.HasPublicationsPlanned {
	case "yes":
		msg.HasPublicationsPlanned = api.Confirmation_CONFIRMATION_YES
	case "no":
		msg.HasPublicationsPlanned = api.Confirmation_CONFIRMATION_NO
	case "dontknow":
		msg.HasPublicationsPlanned = api.Confirmation_CONFIRMATION_DONT_KNOW
	}
	switch p.HasPublishedMaterial {
	case "yes":
		msg.HasPublishedMaterial = api.Confirmation_CONFIRMATION_YES
	case "no":
		msg.HasPublishedMaterial = api.Confirmation_CONFIRMATION_NO
	case "dontknow":
		msg.HasPublishedMaterial = api.Confirmation_CONFIRMATION_DONT_KNOW
	}

	msg.Isbn = p.ISBN

	msg.Issn = p.ISSN

	msg.Issue = p.Issue

	msg.IssueTitle = p.IssueTitle

	msg.Keyword = p.Keyword

	msg.Language = p.Language

	for _, val := range p.LaySummary {
		msg.LaySummary = append(msg.LaySummary, &api.Text{
			Id:   val.ID,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	msg.Legacy = p.Legacy

	for _, val := range p.Supervisor {
		var depts []*api.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, &api.ContributorDepartment{
				Id:   dept.ID,
				Name: dept.Name,
			})
		}
		msg.Supervisor = append(msg.Supervisor, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
			CreditRole: val.CreditRole,
		})
	}

	msg.Volume = p.Volume

	msg.WosId = p.WOSID

	msg.WosType = p.WOSType

	msg.Year = p.Year

	msg.ReportNumber = p.ReportNumber

	msg.ResearchField = p.ResearchField

	msg.ReviewerNote = p.ReviewerNote

	msg.ReviewerTags = p.ReviewerTags

	msg.SeriesTitle = p.SeriesTitle

	msg.SnapshotId = p.SnapshotID

	msg.SourceDb = p.SourceDB

	msg.SourceId = p.SourceID

	msg.SourceRecord = p.SourceRecord

	msg.Locked = p.Locked

	msg.Message = p.Message

	msg.PageCount = p.PageCount
	msg.PageFirst = p.PageFirst
	msg.PageLast = p.PageLast

	msg.PlaceOfPublication = p.PlaceOfPublication

	msg.Publication = p.Publication

	msg.PublicationAbbreviation = p.PublicationAbbreviation

	msg.Publisher = p.Publisher

	msg.PubmedId = p.PubMedID

	msg.VabbId = p.VABBID
	msg.VabbApproved = p.VABBApproved
	msg.VabbType = p.VABBType
	msg.VabbYear = p.VABBYear

	switch p.JournalArticleType {
	case "original":
		msg.JournalArticleType = api.Publication_JOURNAL_ARTICLE_TYPE_ORIGINAL
	case "review":
		msg.JournalArticleType = api.Publication_JOURNAL_ARTICLE_TYPE_REVIEW
	case "letterNote":
		msg.JournalArticleType = api.Publication_JOURNAL_ARTICLE_TYPE_LETTER_NOTE
	case "proceedingsPaper":
		msg.JournalArticleType = api.Publication_JOURNAL_ARTICLE_TYPE_PROCEEDINGS_PAPER
	}

	switch p.ConferenceType {
	case "proceedingsPaper":
		msg.ConferenceType = api.Publication_CONFERENCE_TYPE_PROCEEDINGS_PAPER
	case "abstract":
		msg.ConferenceType = api.Publication_CONFERENCE_TYPE_ABSTRACT
	case "poster":
		msg.ConferenceType = api.Publication_CONFERENCE_TYPE_POSTER
	case "other":
		msg.ConferenceType = api.Publication_CONFERENCE_TYPE_OTHER
	}

	switch p.MiscellaneousType {
	case "artReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_ART_REVIEW
	case "artisticWork":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_ARTISTIC_WORK
	case "bibliography":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_BIBLIOGRAPHY
	case "biography":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_BIOGRAPHY
	case "blogPost":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_BLOG_POST
	case "bookReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_BOOK_REVIEW
	case "correction":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_CORRECTION
	case "dictionaryEntry":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_DICTIONARY_ENTRY
	case "editorialMaterial":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_EDITIORIAL_MATERIAL
	case "encyclopediaEntry":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_ENCYCLOPEDIA_ENTRY
	case "exhibitionReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_EXHIBITION_REVIEW
	case "filmReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_FILM_REVIEW
	case "lectureSpeech":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_LECTURE_SPEECH
	case "lemma":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_LEMMA
	case "magazinePiece":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_MAGAZINE_PIECE
	case "manual":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_MANUAL
	case "musicEdition":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_MUSIC_EDITION
	case "musicReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_MUSIC_REVIEW
	case "newsArticle":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_NEWS_ARTICLE
	case "newspaperPiece":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_NEWSPAPER_PIECE
	case "other":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_OTHER
	case "preprint":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_PREPRINT
	case "product_review":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_PRODUCT_REVIEW
	case "report":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_REPORT
	case "technicalStandard":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_TECHNICAL_STANDARD
	case "textEdition":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_TEXT_EDITION
	case "textTranslation":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_TEXT_TRANSLATION
	case "theatreReview":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_THEATRE_REVIEW
	case "workingPaper":
		msg.MiscellaneousType = api.Publication_MISCELLANEOUS_TYPE_WORKING_PAPER
	}

	for _, val := range p.Project {
		msg.Project = append(msg.Project, &api.RelatedProject{
			Id:   val.ID,
			Name: val.Name,
		})
	}

	for _, val := range p.RelatedDataset {
		msg.RelatedDataset = append(msg.RelatedDataset, &api.RelatedDataset{
			Id: val.ID,
		})
	}

	for _, val := range p.Link {
		l := &api.Link{
			Id:          val.ID,
			Url:         val.URL,
			Description: val.Description,
		}

		switch val.Relation {
		case "data_management_plan":
			l.Relation = api.Link_RELATION_DATA_MANAGEMENT_PLAN
		case "home_page":
			l.Relation = api.Link_RELATION_HOME_PAGE
		case "peer_review_report":
			l.Relation = api.Link_RELATION_PEER_REVIEW_REPORT
		case "related_information":
			l.Relation = api.Link_RELATION_RELATED_INFORMATION
		case "software":
			l.Relation = api.Link_RELATION_SOFTWARE
		case "table_of_contents":
			l.Relation = api.Link_RELATION_TABLE_OF_CONTENTS
		case "main_file":
			l.Relation = api.Link_RELATION_MAIN_FILE
		}

		msg.Link = append(msg.Link, l)
	}

	for _, val := range p.ORCIDWork {
		msg.OrcidWork = append(msg.OrcidWork, &api.OrcidWork{
			Orcid:   val.ORCID,
			PutCode: int32(val.PutCode),
		})
	}

	return msg
}

func MessageToPublication(msg *api.Publication) *models.Publication {
	p := &models.Publication{}

	p.ID = msg.Id

	switch msg.PublicationStatus {
	case api.Publication_PUBLICATION_STATUS_ACCEPTED:
		p.PublicationStatus = "accepted"
	case api.Publication_PUBLICATION_STATUS_UNPUBLISHED:
		p.PublicationStatus = "unpublished"
	case api.Publication_PUBLICATION_STATUS_PUBLISHED:
		p.PublicationStatus = "published"
	}

	switch msg.Type {
	case api.Publication_TYPE_JOURNAL_ARTICLE:
		p.Type = "journal_article"
	case api.Publication_TYPE_BOOK:
		p.Type = "book"
	case api.Publication_TYPE_BOOK_CHAPTER:
		p.Type = "book_chapter"
	case api.Publication_TYPE_BOOK_EDITOR:
		p.Type = "book_editor"
	case api.Publication_TYPE_ISSUE_EDITOR:
		p.Type = "issue_editor"
	case api.Publication_TYPE_CONFERENCE:
		p.Type = "conference"
	case api.Publication_TYPE_DISSERTATION:
		p.Type = "dissertation"
	case api.Publication_TYPE_MISCELLANEOUS:
		p.Type = "miscellaneous"
	}

	switch msg.Status {
	case api.Publication_STATUS_PRIVATE:
		p.Status = "private"
	case api.Publication_STATUS_PUBLIC:
		p.Status = "public"
	case api.Publication_STATUS_DELETED:
		p.Status = "deleted"
	case api.Publication_STATUS_RETURNED:
		p.Status = "returned"
	}

	p.HasBeenPublic = msg.HasBeenPublic

	for _, val := range msg.Abstract {
		p.Abstract = append(p.Abstract, models.Text{
			ID:   val.Id,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	p.AdditionalInfo = msg.AdditionalInfo

	p.AlternativeTitle = msg.AlternativeTitle

	p.ArticleNumber = msg.ArticleNumber

	p.ArxivID = msg.ArxivId

	for _, val := range msg.Author {
		var depts []models.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, models.ContributorDepartment{
				ID:   dept.Id,
				Name: dept.Name,
			})
		}

		p.Author = append(p.Author, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	p.BatchID = msg.BatchId

	switch msg.Classification {
	case api.Publication_CLASSIFICATION_U:
		p.Classification = "U"
	case api.Publication_CLASSIFICATION_A1:
		p.Classification = "A1"
	case api.Publication_CLASSIFICATION_A2:
		p.Classification = "A2"
	case api.Publication_CLASSIFICATION_A3:
		p.Classification = "A3"
	case api.Publication_CLASSIFICATION_A4:
		p.Classification = "A4"
	case api.Publication_CLASSIFICATION_B1:
		p.Classification = "B1"
	case api.Publication_CLASSIFICATION_B2:
		p.Classification = "B2"
	case api.Publication_CLASSIFICATION_B3:
		p.Classification = "B3"
	case api.Publication_CLASSIFICATION_C1:
		p.Classification = "C1"
	case api.Publication_CLASSIFICATION_C3:
		p.Classification = "C3"
	case api.Publication_CLASSIFICATION_D1:
		p.Classification = "D1"
	case api.Publication_CLASSIFICATION_D2:
		p.Classification = "D2"
	case api.Publication_CLASSIFICATION_P1:
		p.Classification = "P1"
	case api.Publication_CLASSIFICATION_V:
		p.Classification = "V"
	}

	if msg.DateCreated != nil {
		t := msg.DateCreated.AsTime()
		p.DateCreated = &t
	}
	if msg.DateUpdated != nil {
		t := msg.DateUpdated.AsTime()
		p.DateUpdated = &t
	}
	if msg.DateFrom != nil {
		t := msg.DateFrom.AsTime()
		p.DateFrom = &t
	}
	if msg.DateUntil != nil {
		t := msg.DateUntil.AsTime()
		p.DateUntil = &t
	}

	p.Extern = msg.Extern

	p.Title = msg.Title

	p.DefensePlace = msg.DefensePlace
	p.DefenseDate = msg.DefenseDate
	p.DefenseTime = msg.DefenseTime

	p.ConferenceName = msg.ConferenceName
	p.ConferenceLocation = msg.ConferenceLocation
	p.ConferenceOrganizer = msg.ConferenceOrganizer
	p.ConferenceStartDate = msg.ConferenceStartDate
	p.ConferenceEndDate = msg.ConferenceEndDate

	for _, val := range msg.Department {
		var depts []models.PublicationDepartmentRef
		for _, dept := range val.Tree {
			depts = append(depts, models.PublicationDepartmentRef{
				ID: dept.Id,
			})
		}

		p.Department = append(p.Department, models.PublicationDepartment{
			ID:   val.Id,
			Tree: depts,
		})
	}

	if msg.Creator != nil {
		p.Creator = &models.PublicationUser{ID: msg.Creator.Id, Name: msg.Creator.Name}
	}

	if msg.User != nil {
		p.User = &models.PublicationUser{ID: msg.User.Id, Name: msg.User.Name}
	}

	if msg.LastUser != nil {
		p.LastUser = &models.PublicationUser{ID: msg.LastUser.Id, Name: msg.LastUser.Name}
	}

	p.DOI = msg.Doi

	p.Edition = msg.Edition

	for _, val := range msg.Editor {
		var depts []models.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, models.ContributorDepartment{
				ID:   dept.Id,
				Name: dept.Name,
			})
		}
		p.Editor = append(p.Editor, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
			CreditRole: val.CreditRole,
		})
	}

	p.EISBN = msg.Eisbn

	p.EISSN = msg.Eissn

	p.ESCIID = msg.EsciId

	for _, val := range msg.File {
		f := &models.PublicationFile{
			ID:           val.Id,
			License:      val.License,
			ContentType:  val.ContentType,
			EmbargoDate:  val.EmbargoDate,
			Name:         val.Name,
			Size:         int(val.Size),
			SHA256:       val.Sha256,
			OtherLicense: val.OtherLicense,
		}

		switch val.AccessLevel {
		case api.File_ACCESS_LEVEL_OPEN_ACCESS:
			f.AccessLevel = "info:eu-repo/semantics/openAccess"
		case api.File_ACCESS_LEVEL_EMBARGOED_ACCESS:
			f.AccessLevel = "info:eu-repo/semantics/embargoedAccess"
		case api.File_ACCESS_LEVEL_RESTRICTED_ACCESS:
			f.AccessLevel = "info:eu-repo/semantics/restrictedAccess"
		case api.File_ACCESS_LEVEL_CLOSED_ACCESS:
			f.AccessLevel = "info:eu-repo/semantics/closedAccess"
		}

		switch val.AccessLevelDuringEmbargo {
		case api.File_ACCESS_LEVEL_CLOSED_ACCESS:
			f.AccessLevelDuringEmbargo = "info:eu-repo/semantics/closedAccess"
		case api.File_ACCESS_LEVEL_RESTRICTED_ACCESS:
			f.AccessLevelDuringEmbargo = "info:eu-repo/semantics/restrictedAccess"
		}

		switch val.AccessLevelAfterEmbargo {
		case api.File_ACCESS_LEVEL_OPEN_ACCESS:
			f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/openAccess"
		case api.File_ACCESS_LEVEL_RESTRICTED_ACCESS:
			f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/restrictedAccess"
		}

		if val.DateCreated != nil {
			t := msg.DateCreated.AsTime()
			f.DateCreated = &t
		}
		if val.DateUpdated != nil {
			t := msg.DateUpdated.AsTime()
			f.DateUpdated = &t
		}

		switch val.PublicationVersion {
		case api.File_PUBLICATION_VERSION_PUBLISHED_VERSION:
			f.PublicationVersion = "publishedVersion"
		case api.File_PUBLICATION_VERSION_AUTHOR_VERSION:
			f.PublicationVersion = "authorVersion"
		case api.File_PUBLICATION_VERSION_ACCEPTED_VERSION:
			f.PublicationVersion = "acceptedVersion"
		case api.File_PUBLICATION_VERSION_UPDATED_VERSION:
			f.PublicationVersion = "updatedVersion"
		}

		switch val.Relation {
		case api.File_RELATION_MAIN_FILE:
			f.Relation = "main_file"
		case api.File_RELATION_COLOPHON:
			f.Relation = "colophon"
		case api.File_RELATION_DATA_FACT_SHEET:
			f.Relation = "data_fact_sheet"
		case api.File_RELATION_PEER_REVIEW_REPORT:
			f.Relation = "peer_review_report"
		case api.File_RELATION_TABLE_OF_CONTENTS:
			f.Relation = "table_of_contents"
		case api.File_RELATION_AGREEMENT:
			f.Relation = "agreement"
		case api.File_RELATION_SUPPLEMENTARY_MATERIAL:
			f.Relation = "supplementary_material"
		}

		p.File = append(p.File, f)
	}

	p.Handle = msg.Handle

	switch msg.HasConfidentialData {
	case api.Confirmation_CONFIRMATION_YES:
		p.HasConfidentialData = "yes"
	case api.Confirmation_CONFIRMATION_NO:
		p.HasConfidentialData = "no"
	case api.Confirmation_CONFIRMATION_DONT_KNOW:
		p.HasConfidentialData = "dontknow"
	}
	switch msg.HasPatentApplication {
	case api.Confirmation_CONFIRMATION_YES:
		p.HasPatentApplication = "yes"
	case api.Confirmation_CONFIRMATION_NO:
		p.HasPatentApplication = "no"
	case api.Confirmation_CONFIRMATION_DONT_KNOW:
		p.HasPatentApplication = "dontknow"
	}
	switch msg.HasPublicationsPlanned {
	case api.Confirmation_CONFIRMATION_YES:
		p.HasPublicationsPlanned = "yes"
	case api.Confirmation_CONFIRMATION_NO:
		p.HasPublicationsPlanned = "no"
	case api.Confirmation_CONFIRMATION_DONT_KNOW:
		p.HasPublicationsPlanned = "dontknow"
	}
	switch msg.HasPublishedMaterial {
	case api.Confirmation_CONFIRMATION_YES:
		p.HasPublishedMaterial = "yes"
	case api.Confirmation_CONFIRMATION_NO:
		p.HasPublishedMaterial = "no"
	case api.Confirmation_CONFIRMATION_DONT_KNOW:
		p.HasPublishedMaterial = "dontknow"
	}

	p.ISBN = msg.Isbn

	p.ISSN = msg.Issn

	p.Issue = msg.Issue

	p.IssueTitle = msg.IssueTitle

	p.Keyword = msg.Keyword

	p.Language = msg.Language

	for _, val := range msg.LaySummary {
		p.LaySummary = append(p.LaySummary, models.Text{
			ID:   val.Id,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range msg.Supervisor {
		var depts []models.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, models.ContributorDepartment{
				ID:   dept.Id,
				Name: dept.Name,
			})
		}
		p.Supervisor = append(p.Supervisor, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
			CreditRole: val.CreditRole,
		})
	}

	p.Legacy = msg.Legacy

	p.Volume = msg.Volume

	p.WOSID = msg.WosId

	p.WOSType = msg.WosType

	p.Year = msg.Year

	p.ReportNumber = msg.ReportNumber

	p.ResearchField = msg.ResearchField

	p.ReviewerNote = msg.ReviewerNote

	p.ReviewerTags = msg.ReviewerTags

	p.SeriesTitle = msg.SeriesTitle

	p.SnapshotID = msg.SnapshotId

	p.SourceDB = msg.SourceDb

	p.SourceID = msg.SourceId

	p.SourceRecord = msg.SourceRecord

	p.Locked = msg.Locked

	p.Message = msg.Message

	p.PageCount = msg.PageCount
	p.PageFirst = msg.PageFirst
	p.PageLast = msg.PageLast

	p.PlaceOfPublication = msg.PlaceOfPublication

	p.Publication = msg.Publication

	p.PublicationAbbreviation = msg.PublicationAbbreviation

	p.Publisher = msg.Publisher

	p.PubMedID = msg.PubmedId

	p.VABBID = msg.VabbId
	p.VABBApproved = msg.VabbApproved
	p.VABBType = msg.VabbType
	p.VABBYear = msg.VabbYear

	switch msg.JournalArticleType {
	case api.Publication_JOURNAL_ARTICLE_TYPE_ORIGINAL:
		p.JournalArticleType = "original"
	case api.Publication_JOURNAL_ARTICLE_TYPE_REVIEW:
		p.JournalArticleType = "review"
	case api.Publication_JOURNAL_ARTICLE_TYPE_LETTER_NOTE:
		p.JournalArticleType = "letterNote"
	case api.Publication_JOURNAL_ARTICLE_TYPE_PROCEEDINGS_PAPER:
		p.JournalArticleType = "proceedingsPaper"
	}

	switch msg.ConferenceType {
	case api.Publication_CONFERENCE_TYPE_PROCEEDINGS_PAPER:
		p.ConferenceType = "proceedingsPaper"
	case api.Publication_CONFERENCE_TYPE_ABSTRACT:
		p.ConferenceType = "abstract"
	case api.Publication_CONFERENCE_TYPE_POSTER:
		p.ConferenceType = "poster"
	case api.Publication_CONFERENCE_TYPE_OTHER:
		p.ConferenceType = "other"
	}

	switch msg.MiscellaneousType {
	case api.Publication_MISCELLANEOUS_TYPE_ART_REVIEW:
		p.MiscellaneousType = "artReview"
	case api.Publication_MISCELLANEOUS_TYPE_ARTISTIC_WORK:
		p.MiscellaneousType = "artisticWork"
	case api.Publication_MISCELLANEOUS_TYPE_BIBLIOGRAPHY:
		p.MiscellaneousType = "bibliography"
	case api.Publication_MISCELLANEOUS_TYPE_BIOGRAPHY:
		p.MiscellaneousType = "biography"
	case api.Publication_MISCELLANEOUS_TYPE_BLOG_POST:
		p.MiscellaneousType = "blogPost"
	case api.Publication_MISCELLANEOUS_TYPE_BOOK_REVIEW:
		p.MiscellaneousType = "bookReview"
	case api.Publication_MISCELLANEOUS_TYPE_CORRECTION:
		p.MiscellaneousType = "correction"
	case api.Publication_MISCELLANEOUS_TYPE_DICTIONARY_ENTRY:
		p.MiscellaneousType = "dictionaryEntry"
	case api.Publication_MISCELLANEOUS_TYPE_EDITIORIAL_MATERIAL:
		p.MiscellaneousType = "editorialMaterial"
	case api.Publication_MISCELLANEOUS_TYPE_ENCYCLOPEDIA_ENTRY:
		p.MiscellaneousType = "encyclopediaEntry"
	case api.Publication_MISCELLANEOUS_TYPE_EXHIBITION_REVIEW:
		p.MiscellaneousType = "exhibitionReview"
	case api.Publication_MISCELLANEOUS_TYPE_FILM_REVIEW:
		p.MiscellaneousType = "filmReview"
	case api.Publication_MISCELLANEOUS_TYPE_LECTURE_SPEECH:
		p.MiscellaneousType = "lectureSpeech"
	case api.Publication_MISCELLANEOUS_TYPE_LEMMA:
		p.MiscellaneousType = "lemma"
	case api.Publication_MISCELLANEOUS_TYPE_MAGAZINE_PIECE:
		p.MiscellaneousType = "magazinePiece"
	case api.Publication_MISCELLANEOUS_TYPE_MANUAL:
		p.MiscellaneousType = "manual"
	case api.Publication_MISCELLANEOUS_TYPE_MUSIC_EDITION:
		p.MiscellaneousType = "musicEdition"
	case api.Publication_MISCELLANEOUS_TYPE_MUSIC_REVIEW:
		p.MiscellaneousType = "musicReview"
	case api.Publication_MISCELLANEOUS_TYPE_NEWS_ARTICLE:
		p.MiscellaneousType = "newsArticle"
	case api.Publication_MISCELLANEOUS_TYPE_NEWSPAPER_PIECE:
		p.MiscellaneousType = "newspaperPiece"
	case api.Publication_MISCELLANEOUS_TYPE_OTHER:
		p.MiscellaneousType = "other"
	case api.Publication_MISCELLANEOUS_TYPE_PREPRINT:
		p.MiscellaneousType = "preprint"
	case api.Publication_MISCELLANEOUS_TYPE_PRODUCT_REVIEW:
		p.MiscellaneousType = "product_review"
	case api.Publication_MISCELLANEOUS_TYPE_REPORT:
		p.MiscellaneousType = "report"
	case api.Publication_MISCELLANEOUS_TYPE_TECHNICAL_STANDARD:
		p.MiscellaneousType = "technicalStandard"
	case api.Publication_MISCELLANEOUS_TYPE_TEXT_EDITION:
		p.MiscellaneousType = "textEdition"
	case api.Publication_MISCELLANEOUS_TYPE_TEXT_TRANSLATION:
		p.MiscellaneousType = "textTranslation"
	case api.Publication_MISCELLANEOUS_TYPE_THEATRE_REVIEW:
		p.MiscellaneousType = "theatreReview"
	case api.Publication_MISCELLANEOUS_TYPE_WORKING_PAPER:
		p.MiscellaneousType = "workingPaper"
	}

	for _, val := range msg.Project {
		// TODO add Name
		p.Project = append(p.Project, models.PublicationProject{
			ID:   val.Id,
			Name: val.Name,
		})
	}

	for _, val := range msg.RelatedDataset {
		p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{
			ID: val.Id,
		})
	}

	for _, val := range msg.Link {
		l := models.PublicationLink{
			ID:          val.Id,
			URL:         val.Url,
			Description: val.Description,
		}

		switch val.Relation {
		case api.Link_RELATION_DATA_MANAGEMENT_PLAN:
			l.Relation = "data_management_plan"
		case api.Link_RELATION_HOME_PAGE:
			l.Relation = "home_page"
		case api.Link_RELATION_PEER_REVIEW_REPORT:
			l.Relation = "peer_review_report"
		case api.Link_RELATION_RELATED_INFORMATION:
			l.Relation = "related_information"
		case api.Link_RELATION_SOFTWARE:
			l.Relation = "software"
		case api.Link_RELATION_TABLE_OF_CONTENTS:
			l.Relation = "table_of_contents"
		case api.Link_RELATION_MAIN_FILE:
			l.Relation = "main_file"
		}

		p.Link = append(p.Link, l)
	}

	for _, val := range msg.OrcidWork {
		p.ORCIDWork = append(p.ORCIDWork, models.PublicationORCIDWork{
			ORCID:   val.Orcid,
			PutCode: int(val.PutCode),
		})
	}

	return p
}
