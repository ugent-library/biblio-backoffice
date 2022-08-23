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
		msg.Type = api.Publication_JOURNAL_ARTICLE
	case "book":
		msg.Type = api.Publication_BOOK
	case "book_chapter":
		msg.Type = api.Publication_BOOK_CHAPTER
	case "book_editor":
		msg.Type = api.Publication_BOOK_EDITOR
	case "issue_editor":
		msg.Type = api.Publication_ISSUE_EDITOR
	case "conference":
		msg.Type = api.Publication_CONFERENCE
	case "dissertation":
		msg.Type = api.Publication_DISSERTATION
	case "miscellaneous":
		msg.Type = api.Publication_MISCELLANEOUS
	}

	switch p.Status {
	case "private":
		msg.Status = api.Publication_PRIVATE
	case "public":
		msg.Status = api.Publication_PUBLIC
	case "deleted":
		msg.Status = api.Publication_DELETED
	case "returned":
		msg.Status = api.Publication_RETURNED
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
		msg.Classification = api.Publication_U
	case "A1":
		msg.Classification = api.Publication_A1
	case "A2":
		msg.Classification = api.Publication_A2
	case "A3":
		msg.Classification = api.Publication_A3
	case "A4":
		msg.Classification = api.Publication_A4
	case "B1":
		msg.Classification = api.Publication_B1
	case "B2":
		msg.Classification = api.Publication_B2
	case "B3":
		msg.Classification = api.Publication_B3
	case "C1":
		msg.Classification = api.Publication_C1
	case "C3":
		msg.Classification = api.Publication_C3
	case "D1":
		msg.Classification = api.Publication_D1
	case "D2":
		msg.Classification = api.Publication_D2
	case "P1":
		msg.Classification = api.Publication_P1
	case "V":
		msg.Classification = api.Publication_V
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
		var tree []*api.Id
		for _, v := range val.Tree {
			tree = append(tree, &api.Id{Id: v.ID})
		}
		msg.Organization = append(msg.Organization, &api.Organization{
			Id:   val.ID,
			Tree: tree,
		})
	}

	msg.CreatorId = p.CreatorID

	msg.UserId = p.UserID

	return msg
}
