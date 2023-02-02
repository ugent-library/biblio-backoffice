package server

import (
	"context"

	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Relate(ctx context.Context, req *api.RelateRequest) (*api.RelateResponse, error) {
	var (
		p   *models.Publication
		d   *models.Dataset
		err error
	)

	switch one := req.One.(type) {
	case *api.RelateRequest_PublicationOne:
		if p, err = s.services.Repository.GetPublication(one.PublicationOne); err != nil {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %s", one.PublicationOne, err)
		}
	case *api.RelateRequest_DatasetOne:
		return nil, status.Error(codes.Internal, "one can only be a publication for now")
	case nil:
		return nil, status.Error(codes.Internal, "one is missing")
	}

	switch two := req.Two.(type) {
	case *api.RelateRequest_PublicationTwo:
		return nil, status.Error(codes.Internal, "two can only be a dataset for now")
	case *api.RelateRequest_DatasetTwo:
		if d, err = s.services.Repository.GetDataset(two.DatasetTwo); err != nil {
			return nil, status.Errorf(codes.Internal, "could not get dataset with id %s: %s", two.DatasetTwo, err)
		}
	case nil:
		return nil, status.Error(codes.Internal, "two is missing")
	}

	if err := s.services.Repository.AddPublicationDataset(p, d, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "could not relate: %s", err)
	}
	if err := s.services.PublicationSearchService.Index(p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to index publication %s, %s", p.ID, err)
	}
	if err := s.services.DatasetSearchService.Index(d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to index dataset %s, %s", d.ID, err)
	}

	return &api.RelateResponse{}, nil
}
