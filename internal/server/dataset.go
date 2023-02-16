package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetDataset(ctx context.Context, req *api.GetDatasetRequest) (*api.GetDatasetResponse, error) {
	d, err := s.services.Repository.GetDataset(req.Id)
	if err != nil {
		// TODO How do we differentiate between errors? e.g. NotFound vs. Internal (database unavailable,...)
		return nil, status.Errorf(codes.Internal, "Could not get dataset with id %s: %s", req.Id, err)
	}

	j, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	apid := &api.Dataset{
		Payload: j,
	}

	res := &api.GetDatasetResponse{Dataset: apid}

	return res, nil
}

func (s *server) GetAllDatasets(req *api.GetAllDatasetsRequest, stream api.Biblio_GetAllDatasetsServer) (err error) {
	// TODO errors in EachDataset aren't caught and pushed upstream. Returning 'false' in the callback
	//   breaks the loop, but EachDataset will return 'nil'.
	ErrorStream := s.services.Repository.EachDataset(func(d *models.Dataset) bool {
		j, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		apid := &api.Dataset{
			Payload: j,
		}

		res := &api.GetAllDatasetsResponse{Dataset: apid}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if ErrorStream != nil {
		return status.Errorf(codes.Internal, "could not get all datasets: %s", ErrorStream)
	}

	return nil
}

func (s *server) SearchDatasets(ctx context.Context, req *api.SearchDatasetsRequest) (*api.SearchDatasetsResponse, error) {
	page := 1
	if req.Limit > 0 {
		page = int(req.Offset)/int(req.Limit) + 1
	}
	args := models.NewSearchArgs().WithQuery(req.Query).WithPage(page)
	hits, err := s.services.DatasetSearchService.Search(args)
	if err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "Could not search datasets: %s :: %s", req.Query, err)
	}

	res := &api.SearchDatasetsResponse{
		Limit:  int32(hits.Limit),
		Offset: int32(hits.Offset),
		Total:  int32(hits.Total),
		Hits:   make([]*api.Dataset, len(hits.Hits)),
	}
	for i, d := range hits.Hits {
		j, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		apid := &api.Dataset{
			Payload: j,
		}
		res.Hits[i] = apid
	}

	return res, nil
}

func (s *server) UpdateDataset(ctx context.Context, req *api.UpdateDatasetRequest) (*api.UpdateDatasetResponse, error) {
	d := &models.Dataset{}
	if err := json.Unmarshal(req.Dataset.Payload, d); err != nil {
		log.Fatal(err)
	}

	// TODO Fetch user information via better authentication (no basic auth)
	user := &models.User{
		ID:       "n/a",
		FullName: "system user",
	}

	if err := d.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for dataset %s: %s", d.ID, err)
	}

	if err := s.services.Repository.UpdateDataset(d.SnapshotID, d, user); err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "failed to store dataset %s, %s", d.ID, err)
	}
	if err := s.services.DatasetSearchService.Index(d); err != nil {
		// TODO How do we differentiate between errors
		return nil, status.Errorf(codes.Internal, "failed to index dataset %s, %s", d.ID, err)
	}

	return &api.UpdateDatasetResponse{}, nil
}

// TODO catch indexing errors
func (s *server) AddDatasets(stream api.Biblio_AddDatasetsServer) error {
	ctx := context.Background()

	bi, err := s.services.DatasetSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			msg := fmt.Errorf("failed to index: %w", err).Error()
			// TODO catch error
			stream.Send(&api.AddDatasetsResponse{Message: msg})
		},
		OnIndexError: func(id string, err error) {
			msg := fmt.Errorf("failed to index publication %s: %w", id, err).Error()
			// TODO catch error
			stream.Send(&api.AddDatasetsResponse{Message: msg})
		},
	})
	if err != nil {
		return err
	}
	defer bi.Close(ctx)

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

		d := &models.Dataset{}
		if err := json.Unmarshal(req.Dataset.Payload, d); err != nil {
			log.Fatal(err)
		}

		if d.ID == "" {
			d.ID = ulid.Make().String()
		}

		if d.Status == "" {
			d.Status = "private"
		}

		for i, val := range d.Abstract {
			if val.ID == "" {
				val.ID = ulid.Make().String()
			}
			d.Abstract[i] = val
		}

		// TODO this should return structured messages (see validate)
		if err := d.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for dataset %s at line %d: %s", d.ID, seq, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.SaveDataset(d, nil); err != nil {
			msg := fmt.Errorf("failed to store dataset %s at line %d: %s", d.ID, seq, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := bi.Index(ctx, d); err != nil {
			msg := fmt.Errorf("failed to index dataset %s at line %d: %w", d.ID, seq, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
	}
}

func (s *server) ImportDatasets(stream api.Biblio_ImportDatasetsServer) error {
	ctx := context.Background()

	bi, err := s.services.DatasetSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			msg := fmt.Errorf("failed to index: %w", err).Error()
			// TODO catch error
			stream.Send(&api.ImportDatasetsResponse{Message: msg})
		},
		OnIndexError: func(id string, err error) {
			msg := fmt.Errorf("failed to index dataset %s: %w", id, err).Error()
			// TODO catch error
			stream.Send(&api.ImportDatasetsResponse{Message: msg})
		},
	})
	if err != nil {
		return err
	}
	defer bi.Close(ctx)
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

		d := &models.Dataset{}
		if err := json.Unmarshal(req.Dataset.Payload, d); err != nil {
			log.Fatal(err)
		}

		// TODO this should return structured messages (see validate)
		if err := d.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for dataset %s at line %d: %s", d.ID, seq, err).Error()
			if err = stream.Send(&api.ImportDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.ImportCurrentDataset(d); err != nil {
			msg := fmt.Errorf("failed to store dataset %s at line %d: %s", d.ID, seq, err).Error()
			if err = stream.Send(&api.ImportDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := bi.Index(ctx, d); err != nil {
			msg := fmt.Errorf("failed to index dataset %s at line %d: %w", d.ID, seq, err).Error()
			if err = stream.Send(&api.ImportDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
	}
}

func (s *server) GetDatasetHistory(req *api.GetDatasetHistoryRequest, stream api.Biblio_GetDatasetHistoryServer) (err error) {
	errorStream := s.services.Repository.DatasetHistory(req.Id, func(d *models.Dataset) bool {
		j, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		apid := &api.Dataset{
			Payload: j,
		}

		res := &api.GetDatasetHistoryResponse{Dataset: apid}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get dataset history: %s", errorStream)
	}

	return nil
}

func (s *server) PurgeDataset(ctx context.Context, req *api.PurgeDatasetRequest) (*api.PurgeDatasetResponse, error) {
	if err := s.services.Repository.PurgeDataset(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset with id %s: %s", req.Id, err)
	}
	if err := s.services.DatasetSearchService.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset from index with id %s: %s", req.Id, err)
	}

	return &api.PurgeDatasetResponse{}, nil
}

func (s *server) PurgeAllDatasets(ctx context.Context, req *api.PurgeAllDatasetsRequest) (*api.PurgeAllDatasetsResponse, error) {
	if err := s.services.Repository.PurgeAllDatasets(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all datasets: %s", err)
	}
	if err := s.services.DatasetSearchService.DeleteAll(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete dataset index: %w", err)
	}
	return &api.PurgeAllDatasetsResponse{}, nil
}

func (s *server) ValidateDatasets(stream api.Biblio_ValidateDatasetsServer) error {
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

		d := &models.Dataset{}
		if err := json.Unmarshal(req.Dataset.Payload, d); err != nil {
			log.Fatal(err)
		}

		err = d.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			res := &api.ValidateDatasetsResponse{Seq: seq, Id: d.ID, Message: validationErrs.Error()}
			if err = stream.Send(res); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}
