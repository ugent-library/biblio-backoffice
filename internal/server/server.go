package server

import (
	"context"

	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"google.golang.org/grpc"
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
	return msg
}
