package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetDataset(ctx context.Context, req *api.GetDatasetRequest) (*api.GetDatasetResponse, error) {
	p, err := s.services.Repo.GetDataset(req.Id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not find dataset with id %s", req.Id).Error())
			return &api.GetDatasetResponse{
				Response: &api.GetDatasetResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get dataset with id %s: %v", req.Id, err)
		}
	}

	j, err := json.Marshal(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal dataset with id %s: %v", req.Id, err)
	}

	res := &api.GetDatasetResponse{
		Response: &api.GetDatasetResponse_Dataset{
			Dataset: &api.Dataset{
				Payload: j,
			},
		},
	}

	return res, nil
}

func (s *server) GetAllDatasets(req *api.GetAllDatasetsRequest, stream api.Biblio_GetAllDatasetsServer) (err error) {
	var callbackErr error

	streamErr := s.services.Repo.EachDataset(func(d *models.Dataset) bool {
		j, err := json.Marshal(d)
		if err != nil {
			grpcError := status.New(codes.Internal, fmt.Errorf("could not marshal dataset with id %s: %v", d.ID, err).Error())
			grpcError, err = grpcError.WithDetails(req)
			if err != nil {
				callbackErr = err
				return false
			}

			stream.Send(&api.GetAllDatasetsResponse{
				Response: &api.GetAllDatasetsResponse_Error{
					Error: grpcError.Proto(),
				},
			})

			return true
		}

		res := &api.GetAllDatasetsResponse{
			Response: &api.GetAllDatasetsResponse_Dataset{
				Dataset: &api.Dataset{
					Payload: j,
				},
			},
		}

		if err = stream.Send(res); err != nil {
			callbackErr = err
			return false
		}

		return true
	})

	if streamErr != nil {
		return status.Errorf(codes.Internal, "could not get all datasets: %v", streamErr)
	}

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "could not get all datasets: %v", callbackErr)
	}

	return nil
}

func (s *server) SearchDatasets(ctx context.Context, req *api.SearchDatasetsRequest) (*api.SearchDatasetsResponse, error) {
	page := 1
	if req.Limit > 0 {
		page = int(req.Offset)/int(req.Limit) + 1
	}
	args := models.NewSearchArgs().WithQuery(req.Query).WithPage(page)
	hits, err := s.services.DatasetSearchIndex.Search(args)
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
		grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
		return &api.UpdateDatasetResponse{
			Response: &api.UpdateDatasetResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	if err := d.Validate(); err != nil {
		grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate dataset %s: %v", d.ID, err).Error())
		return &api.UpdateDatasetResponse{
			Response: &api.UpdateDatasetResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	err := s.services.Repo.UpdateDataset(d.SnapshotID, d, nil)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		grpcErr := status.New(codes.Internal, fmt.Errorf("failed to update dataset: conflict detected for dataset[snapshot_id: %s, id: %s] : %v", d.SnapshotID, d.ID, err).Error())
		return &api.UpdateDatasetResponse{
			Response: &api.UpdateDatasetResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update dataset[snapshot_id: %s, id: %s], %s", d.SnapshotID, d.ID, err)
	}

	return &api.UpdateDatasetResponse{}, nil
}

func (s *server) AddDatasets(stream api.Biblio_AddDatasetsServer) error {
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
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.AddDatasetsResponse{
				Response: &api.AddDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add datasets: %v", err)
			}
			continue
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

		if err := d.Validate(); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate dataset %s at line %d: %v", d.ID, seq, err).Error())
			if err = stream.Send(&api.AddDatasetsResponse{
				Response: &api.AddDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add datasets: %v", err)
			}
			continue
		}

		if err := s.services.Repo.SaveDataset(d, nil); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to store dataset %s at line %d: %s", d.ID, seq, err).Error())
			if err = stream.Send(&api.AddDatasetsResponse{
				Response: &api.AddDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add datasets: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.AddDatasetsResponse{
			Response: &api.AddDatasetsResponse_Message{
				Message: fmt.Sprintf("stored and indexed dataset %s at line %d", d.ID, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to add datasets: %v", err)
		}
	}
}

func (s *server) ImportDatasets(stream api.Biblio_ImportDatasetsServer) error {
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
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.ImportDatasetsResponse{
				Response: &api.ImportDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import datasets: %v", err)
			}
			continue
		}

		if err := d.Validate(); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate dataset at line %d: %v", seq, err).Error())
			if err = stream.Send(&api.ImportDatasetsResponse{
				Response: &api.ImportDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import datasets: %v", err)
			}
			continue
		}

		if err := s.services.Repo.ImportDataset(d); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to store dataset %s at line %d: %s", d.ID, seq, err).Error())
			if err = stream.Send(&api.ImportDatasetsResponse{
				Response: &api.ImportDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import datasets: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.ImportDatasetsResponse{
			Response: &api.ImportDatasetsResponse_Message{
				Message: fmt.Sprintf("stored and indexed dataset %s at line %d", d.ID, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to import datasets: %v", err)
		}
	}
}

func (s *server) MutateDatasets(stream api.Biblio_MutateDatasetsServer) error {
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

		mut := repositories.Mutation{
			Op:   req.Op,
			Args: req.Args,
		}

		if err := s.services.Repo.MutateDataset(req.Id, nil, mut); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to mutate dataset %s at line %d: %s", req.Id, seq, err).Error())
			if err = stream.Send(&api.MutateResponse{
				Response: &api.MutateResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to mutate dataset: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.MutateResponse{
			Response: &api.MutateResponse_Message{
				Message: fmt.Sprintf("mutated dataset %s at line %d", req.Id, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to mutate dataset: %v", err)
		}
	}
}

func (s *server) GetDatasetHistory(req *api.GetDatasetHistoryRequest, stream api.Biblio_GetDatasetHistoryServer) (err error) {
	var callbackErr error
	streamErr := s.services.Repo.DatasetHistory(req.Id, func(d *models.Dataset) bool {
		j, err := json.Marshal(d)
		if err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.GetDatasetHistoryResponse{
				Response: &api.GetDatasetHistoryResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				callbackErr = err
				return false
			}
		} else {
			if err = stream.Send(&api.GetDatasetHistoryResponse{
				Response: &api.GetDatasetHistoryResponse_Dataset{
					Dataset: &api.Dataset{
						Payload: j,
					},
				},
			}); err != nil {
				// PICK UP ERROR
				callbackErr = err
				return false
			}
		}

		return true
	})

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "could not get dataset history: %v", callbackErr)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "could not get dataset history: %v	", streamErr)
	}

	return nil
}

func (s *server) PurgeDataset(ctx context.Context, req *api.PurgeDatasetRequest) (*api.PurgeDatasetResponse, error) {
	_, err := s.services.Repo.GetDataset(req.Id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			grpcErr := status.New(codes.NotFound, fmt.Errorf("could not find dataset with id %s", req.Id).Error())
			return &api.PurgeDatasetResponse{
				Response: &api.PurgeDatasetResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get dataset with id %s: %v", req.Id, err)
		}
	}

	// TODO purgeDataset doesn't return an error if the record for req.Id can't be found
	if err := s.services.Repo.PurgeDataset(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset with id %s: %s", req.Id, err)
	}

	// TODO this will complain if the above didn't throw a 'not found' error
	if err := s.services.DatasetSearchIndex.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset from index with id %s: %s", req.Id, err)
	}

	return &api.PurgeDatasetResponse{
		Response: &api.PurgeDatasetResponse_Ok{
			Ok: true,
		},
	}, nil
}

func (s *server) PurgeAllDatasets(ctx context.Context, req *api.PurgeAllDatasetsRequest) (*api.PurgeAllDatasetsResponse, error) {
	if !req.Confirm {
		grpcErr := status.New(codes.Internal, fmt.Errorf("confirm property in request is not set to true").Error())
		return &api.PurgeAllDatasetsResponse{
			Response: &api.PurgeAllDatasetsResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	if err := s.services.Repo.PurgeAllDatasets(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all datasets: %s", err)
	}

	if err := s.services.DatasetSearchIndex.DeleteAll(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete dataset from index: %s", err)
	}

	return &api.PurgeAllDatasetsResponse{
		Response: &api.PurgeAllDatasetsResponse_Ok{
			Ok: true,
		},
	}, nil
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
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.ValidateDatasetsResponse{
				Response: &api.ValidateDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to validate datasets: %v", err)
			}
			continue
		}

		err = d.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			if err = stream.Send(&api.ValidateDatasetsResponse{
				Response: &api.ValidateDatasetsResponse_Results{
					Results: &api.ValidateResults{
						Seq:     seq,
						Id:      d.ID,
						Message: validationErrs.Error(),
					},
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to validate datasets: %v", err)
			}
		} else if err != nil {
			return status.Errorf(codes.Internal, "failed to validate datasets: %v", err)
		}
	}
}

func (s *server) ReindexDatasets(req *api.ReindexDatasetsRequest, stream api.Biblio_ReindexDatasetsServer) error {
	startTime := time.Now()
	indexed := 0
	reported := 0

	if err := stream.Send(&api.ReindexDatasetsResponse{
		Response: &api.ReindexDatasetsResponse_Message{
			Message: "Indexing to a new index",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	var swErr error
	var swIdxErr error

	switcher, err := s.services.SearchService.NewDatasetIndexSwitcher(backends.BulkIndexerConfig{
		OnError: func(err error) {
			grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index dataset: %v", err).Error())
			if err = stream.Send(&api.ReindexDatasetsResponse{
				Response: &api.ReindexDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				swErr = err
			}
		},
		OnIndexError: func(id string, err error) {
			grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index dataset %s: %w", id, err).Error())
			if err = stream.Send(&api.ReindexDatasetsResponse{
				Response: &api.ReindexDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				swIdxErr = err
			}
		},
	})

	if err != nil {
		return status.Errorf(codes.Internal, "failed to start an indexer: %s", err)
	}

	ctx := stream.Context()
	var callbackErr error

	streamErr := s.services.Repo.EachDataset(func(d *models.Dataset) bool {
		if err := switcher.Index(ctx, d); err != nil {
			grpcErr := status.New(codes.Internal, fmt.Errorf("indexing failed for dataset [id: %s] : %s", d.ID, err).Error())
			if err = stream.Send(&api.ReindexDatasetsResponse{
				Response: &api.ReindexDatasetsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				// PICK UP ERROR
				return false
			}
		}

		indexed++

		// progress message
		if indexed-500 == reported {
			stream.Send(&api.ReindexDatasetsResponse{
				Response: &api.ReindexDatasetsResponse_Message{
					Message: fmt.Sprintf("Indexing %d datasets...", indexed),
				},
			})
			reported = indexed
		}

		return true
	})

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	if swErr != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", swErr)
	}

	if swIdxErr != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", swIdxErr)
	}

	if err := stream.Send(&api.ReindexDatasetsResponse{
		Response: &api.ReindexDatasetsResponse_Message{
			Message: fmt.Sprintf("Indexed %d datasets", indexed),
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	if err := stream.Send(&api.ReindexDatasetsResponse{
		Response: &api.ReindexDatasetsResponse_Message{
			Message: "Switching to new index...",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	if err := switcher.Switch(ctx); err != nil {
		return status.Errorf(codes.Internal, "indexing failed: %v", err)
	}

	endTime := time.Now()

	if err := stream.Send(&api.ReindexDatasetsResponse{
		Response: &api.ReindexDatasetsResponse_Message{
			Message: "Indexing changes since start of reindex",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	for {
		indexed = 0

		var biErr error
		var biIdxErr error

		bi, err := s.services.SearchService.NewDatasetBulkIndexer(backends.BulkIndexerConfig{
			OnError: func(err error) {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index dataset: %v", err).Error())
				if err = stream.Send(&api.ReindexDatasetsResponse{
					Response: &api.ReindexDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					biErr = err
				}
			},
			OnIndexError: func(id string, err error) {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index dataset %s: %w", id, err).Error())
				if err = stream.Send(&api.ReindexDatasetsResponse{
					Response: &api.ReindexDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					biErr = err
				}
			},
		})

		if err != nil {
			return status.Errorf(codes.Internal, "failed to start an indexer: %s", err)
		}

		defer bi.Close(ctx)

		var callbackErr error
		streamErr := s.services.Repo.DatasetsBetween(startTime, endTime, func(d *models.Dataset) bool {
			if err := bi.Index(ctx, d); err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("indexing failed for dataset [id: %s] : %s", d.ID, err).Error())
				if err = stream.Send(&api.ReindexDatasetsResponse{
					Response: &api.ReindexDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}
			}

			indexed++

			return true
		})

		if callbackErr != nil {
			return status.Errorf(codes.Internal, "failed to index datasets: %v", callbackErr)
		}

		if streamErr != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", streamErr)
		}

		if err != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", err)
		}

		if biErr != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", biErr)
		}

		if biIdxErr != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", biErr)
		}

		if indexed == 0 {
			break
		}

		if err := stream.Send(&api.ReindexDatasetsResponse{
			Response: &api.ReindexDatasetsResponse_Message{
				Message: fmt.Sprintf("Indexed %d datasets", indexed),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
		}

		startTime = endTime
		endTime = time.Now()
	}

	if err := stream.Send(&api.ReindexDatasetsResponse{
		Response: &api.ReindexDatasetsResponse_Message{
			Message: "Done",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index datasets: %v", err)
	}

	return nil
}

func (s *server) CleanupDatasets(req *api.CleanupDatasetsRequest, stream api.Biblio_CleanupDatasetsServer) error {
	var callbackErr error

	count := 0
	streamErr := s.services.Repo.EachDataset(func(d *models.Dataset) bool {
		// guard
		fixed := false

		// correctly set HasBeenPublic (only needs to run once)
		if d.Status == "deleted" && !d.HasBeenPublic {
			s.services.Repo.PublicationHistory(d.ID, func(dd *models.Publication) bool {
				if dd.Status == "public" {
					d.HasBeenPublic = true
					fixed = true
					return false
				}
				return true
			})
		}

		// remove empty strings from string array
		vacuumArray := func(old_values []string) []string {
			var newVals []string
			for _, val := range old_values {
				newVal := strings.TrimSpace(val)
				if newVal != "" {
					newVals = append(newVals, val)
				}
				if val != newVal || newVal == "" {
					fixed = true
				}
			}
			return newVals
		}

		d.Format = vacuumArray(d.Format)
		d.Keyword = vacuumArray(d.Keyword)
		d.Language = vacuumArray(d.Language)
		d.ReviewerTags = vacuumArray(d.ReviewerTags)
		for _, author := range d.Author {
			author.CreditRole = vacuumArray(author.CreditRole)
		}
		for _, contributor := range d.Contributor {
			contributor.CreditRole = vacuumArray(contributor.CreditRole)
		}

		// save record if changed
		if fixed {
			d.UserID = ""
			d.User = nil

			if err := d.Validate(); err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to validate dataset[snapshot_id: %s, id: %s] : %v", d.SnapshotID, d.ID, err).Error())
				if err = stream.Send(&api.CleanupDatasetsResponse{
					Response: &api.CleanupDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				return true
			}

			err := s.services.Repo.UpdateDataset(d.SnapshotID, d, nil)
			if err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to update dataset[snapshot_id: %s, id: %s] : %v", d.SnapshotID, d.ID, err).Error())
				if err = stream.Send(&api.CleanupDatasetsResponse{
					Response: &api.CleanupDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				return true
			}

			var conflict *snapstore.Conflict
			if errors.As(err, &conflict) {
				grpcErr := status.New(codes.Internal, fmt.Errorf("conflict detected for dataset[snapshot_id: %s, id: %s] : %v", d.SnapshotID, d.ID, err).Error())
				if err = stream.Send(&api.CleanupDatasetsResponse{
					Response: &api.CleanupDatasetsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				return true
			}

			if err := stream.Send(&api.CleanupDatasetsResponse{
				Response: &api.CleanupDatasetsResponse_Message{
					Message: fmt.Sprintf("fixed dataset[snapshot_id: %s, id: %s]", d.SnapshotID, d.ID),
				},
			}); err != nil {
				callbackErr = err
				return false
			}

			count += 1
		}

		return true
	})

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "failed to complete cleanup: %v", callbackErr)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "could not complete cleanup: %v", streamErr)
	}

	if err := stream.Send(&api.CleanupDatasetsResponse{
		Response: &api.CleanupDatasetsResponse_Message{
			Message: fmt.Sprintf("done. cleaned %d datasets.", count),
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to clean up datasets: %v", err)
	}

	return nil
}
