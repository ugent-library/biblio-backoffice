package server

import (
	"context"

	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	api.UnimplementedBiblioServer
	services *backends.Services
}

func New(services *backends.Services) *grpc.Server {
	gsrv := grpc.NewServer()
	srv := &server{services: services}
	api.RegisterBiblioServer(gsrv, srv)
	return gsrv
}

func (s *server) GetPublication(ctx context.Context, req *api.GetPublicationRequest) (*api.GetPublicationResponse, error) {
	pub, err := s.services.Repository.GetPublication(req.Id)
	if err != nil {
		return nil, err
	}

	res := &api.GetPublicationResponse{Publication: publicationToMessage(pub)}

	return res, nil
}

// TODO keep checking if new publications become available? (see Distributed services with Go)
func (s *server) GetAllPublications(req *api.GetAllPublicationsRequest, stream api.Biblio_GetAllPublicationsServer) error {
	c, err := s.services.Repository.GetAllPublications()
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			if !c.HasNext() {
				return nil
			}
			snap, err := c.Next()
			if err != nil {
				return err
			}
			p := &models.Publication{}
			if err := snap.Scan(p); err != nil {
				return err
			}
			p.SnapshotID = snap.SnapshotID
			p.DateFrom = snap.DateFrom
			p.DateUntil = snap.DateUntil

			res := &api.GetAllPublicationsResponse{Publication: publicationToMessage(p)}

			if err = stream.Send(res); err != nil {
				return err
			}
		}
	}
}

func publicationToMessage(p *models.Publication) *api.Publication {
	msg := &api.Publication{}

	msg.Id = p.ID

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

	for _, val := range p.Abstract {
		msg.Abstract = append(msg.Abstract, &api.Text{
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	msg.AdditionalInfo = p.AdditionalInfo

	msg.AlternativeTitle = p.AlternativeTitle

	msg.ArticleNumber = p.ArticleNumber

	msg.ArxivId = p.ArxivID

	for _, val := range p.Author {
		msg.Author = append(msg.Author, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
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
		msg.Organization = append(msg.Organization, &api.RelatedOrganization{
			Id: val.ID,
		})
	}

	msg.CreatorId = p.CreatorID

	msg.UserId = p.UserID

	msg.Doi = p.DOI

	msg.Edition = p.Edition

	for _, val := range p.Editor {
		msg.Editor = append(msg.Editor, &api.Contributor{
			Id:        val.ID,
			Orcid:     val.ORCID,
			LocalId:   val.UGentID,
			FirstName: val.FirstName,
			LastName:  val.LastName,
			FullName:  val.FullName,
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
			Embargo:      val.Embargo,
			Name:         val.Name,
			Size:         int32(val.Size),
			Sha256:       val.SHA256,
			OtherLicense: val.OtherLicense,
		}

		switch val.AccessLevel {
		case "open_access":
			f.AccessLevel = api.File_ACCESS_LEVEL_OPEN_ACCESS
		case "local":
			f.AccessLevel = api.File_ACCESS_LEVEL_LOCAL
		case "closed":
			f.AccessLevel = api.File_ACCESS_LEVEL_CLOSED
		}

		if val.DateCreated != nil {
			f.DateCreated = timestamppb.New(*val.DateCreated)
		}
		if val.DateUpdated != nil {
			f.DateUpdated = timestamppb.New(*val.DateUpdated)
		}

		switch val.EmbargoTo {
		case "open_access":
			f.EmbargoTo = api.File_ACCESS_LEVEL_OPEN_ACCESS
		case "local":
			f.EmbargoTo = api.File_ACCESS_LEVEL_LOCAL
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
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range p.Supervisor {
		msg.Supervisor = append(msg.Supervisor, &api.Contributor{
			Id:        val.ID,
			Orcid:     val.ORCID,
			LocalId:   val.UGentID,
			FirstName: val.FirstName,
			LastName:  val.LastName,
			FullName:  val.FullName,
		})
	}

	msg.Url = p.URL

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

	msg.SourceId = p.SnapshotID

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
			Id: val.ID,
		})
	}

	for _, val := range p.RelatedDataset {
		msg.Dataset = append(msg.Dataset, &api.RelatedDataset{
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
