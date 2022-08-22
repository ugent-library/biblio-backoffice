package server

import (
	"context"

	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/backends"
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

	res := &api.GetPublicationResponse{Publication: &api.Publication{}}
	res.Publication.Id = pub.ID

	return res, nil
}
