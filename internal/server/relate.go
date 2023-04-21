package server

import (
	"context"
	"errors"
	"fmt"

	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
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
		p, err = s.services.Repository.GetPublication(one.PublicationOne)
		if err != nil {
			if errors.Is(err, backends.ErrNotFound) {
				grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not find publication with id %s", one.PublicationOne).Error())
				return &api.RelateResponse{
					Response: &api.RelateResponse_Error{
						Error: grpcErr.Proto(),
					},
				}, nil
			} else {
				return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %v", one.PublicationOne, err)
			}
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
		d, err = s.services.Repository.GetDataset(two.DatasetTwo)
		if err != nil {
			if errors.Is(err, backends.ErrNotFound) {
				grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not find dataset with id %s", two.DatasetTwo).Error())
				return &api.RelateResponse{
					Response: &api.RelateResponse_Error{
						Error: grpcErr.Proto(),
					},
				}, nil
			} else {
				return nil, status.Errorf(codes.Internal, "could not get dataset with id %s: %v", two.DatasetTwo, err)
			}
		}
	case nil:
		return nil, status.Error(codes.Internal, "two is missing")
	}

	err = s.services.Repository.AddPublicationDataset(p, d, nil)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		grpcErr := status.New(codes.Internal, fmt.Errorf("could not relate publication and dataset: conflict detected").Error())
		return &api.RelateResponse{
			Response: &api.RelateResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update dataset[snapshot_id: %s, id: %s], %s", d.SnapshotID, d.ID, err)
	}

	if err := s.services.PublicationSearchService.Index(p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to index publication %s, %s", p.ID, err)
	}

	if err := s.services.DatasetSearchService.Index(d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to index dataset %s, %s", d.ID, err)
	}

	return &api.RelateResponse{
		Response: &api.RelateResponse_Message{
			Message: fmt.Sprintf("related: publication[id: %s] -> dataset[id: %s]", p.ID, d.ID),
		},
	}, nil
}
