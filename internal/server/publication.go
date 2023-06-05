package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetPublication(ctx context.Context, req *api.GetPublicationRequest) (*api.GetPublicationResponse, error) {
	p, err := s.services.Repository.GetPublication(req.Id)

	if err != nil {
		if errors.Is(err, backends.ErrNotFound) {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not find publication with id %s", req.Id).Error())
			return &api.GetPublicationResponse{
				Response: &api.GetPublicationResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %v", err)
		}
	}

	j, err := json.Marshal(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal publication with id %s: %v", err)
	}

	res := &api.GetPublicationResponse{
		Response: &api.GetPublicationResponse_Publication{
			Publication: &api.Publication{
				Payload: j,
			},
		},
	}

	return res, nil
}

func (s *server) GetAllPublications(req *api.GetAllPublicationsRequest, stream api.Biblio_GetAllPublicationsServer) (err error) {
	var callbackErr error

	streamErr := s.services.Repository.EachPublication(func(p *models.Publication) bool {
		j, err := json.Marshal(p)
		if err != nil {
			grpcError := status.New(codes.Internal, fmt.Errorf("could not marshal publication with id %s: %v", p.ID, err).Error())
			grpcError, err = grpcError.WithDetails(req)
			if err != nil {
				callbackErr = err
				return false
			}

			stream.Send(&api.GetAllPublicationsResponse{
				Response: &api.GetAllPublicationsResponse_Error{
					Error: grpcError.Proto(),
				},
			})

			return true
		}

		res := &api.GetAllPublicationsResponse{
			Response: &api.GetAllPublicationsResponse_Publication{
				Publication: &api.Publication{
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
		return status.Errorf(codes.Internal, "could not get all publications: %v", streamErr)
	}

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "could not get all publications: %v", callbackErr)
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
		return nil, status.Errorf(codes.Internal, "could not search publications: %s :: %s", req.Query, err)
	}

	res := &api.SearchPublicationsResponse{
		Limit:  int32(hits.Limit),
		Offset: int32(hits.Offset),
		Total:  int32(hits.Total),
		Hits:   make([]*api.Publication, len(hits.Hits)),
	}
	for i, p := range hits.Hits {
		j, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}
		apip := &api.Publication{
			Payload: j,
		}

		res.Hits[i] = apip
	}

	return res, nil
}

func (s *server) UpdatePublication(ctx context.Context, req *api.UpdatePublicationRequest) (*api.UpdatePublicationResponse, error) {
	p := &models.Publication{}
	if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
		grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
		return &api.UpdatePublicationResponse{
			Response: &api.UpdatePublicationResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	if err := p.Validate(); err != nil {
		grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate publication %s: %v", p.ID, err).Error())
		return &api.UpdatePublicationResponse{
			Response: &api.UpdatePublicationResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	err := s.services.Repository.UpdatePublication(p.SnapshotID, p, nil)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		grpcErr := status.New(codes.Internal, fmt.Errorf("failed to update publication: conflict detected for publication[snapshot_id: %s, id: %s] : %v", p.SnapshotID, p.ID, err).Error())
		return &api.UpdatePublicationResponse{
			Response: &api.UpdatePublicationResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update publication[snapshot_id: %s, id: %s], %s", p.SnapshotID, p.ID, err)
	}

	if err := s.services.PublicationSearchService.Index(p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update publication[snapshot_id: %s, id: %s], %s", p.SnapshotID, p.ID, err)
	}

	return &api.UpdatePublicationResponse{}, nil
}

func (s *server) AddPublications(stream api.Biblio_AddPublicationsServer) error {
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

		p := &models.Publication{}
		if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.AddPublicationsResponse{
				Response: &api.AddPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add publications: %v", err)
			}
			continue
		}

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

		if err := p.Validate(); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate publication %s at line %d: %v", p.ID, seq, err).Error())
			if err = stream.Send(&api.AddPublicationsResponse{
				Response: &api.AddPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add publications: %v", err)
			}
			continue
		}

		if err := s.services.Repository.SavePublication(p, nil); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error())
			if err = stream.Send(&api.AddPublicationsResponse{
				Response: &api.AddPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to add publications: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.AddPublicationsResponse{
			Response: &api.AddPublicationsResponse_Message{
				Message: fmt.Sprintf("stored and indexed publication %s at line %d", p.ID, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to add publications: %v", err)
		}
	}
}

func (s *server) ImportPublications(stream api.Biblio_ImportPublicationsServer) error {
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

		p := &models.Publication{}
		if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.ImportPublicationsResponse{
				Response: &api.ImportPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import publications: %v", err)
			}
			continue
		}

		if err := p.Validate(); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to validate publication at line %d: %v", seq, err).Error())
			if err = stream.Send(&api.ImportPublicationsResponse{
				Response: &api.ImportPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import publications: %v", err)
			}
			continue
		}

		if err := s.services.Repository.ImportPublication(p); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error())
			if err = stream.Send(&api.ImportPublicationsResponse{
				Response: &api.ImportPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to import publications: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.ImportPublicationsResponse{
			Response: &api.ImportPublicationsResponse_Message{
				Message: fmt.Sprintf("stored and indexed publication %s at line %d", p.ID, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to import publications: %v", err)
		}
	}
}

func (s *server) MutatePublications(stream api.Biblio_MutatePublicationsServer) error {
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

		mut := backends.Mutation{
			Op:   req.Op,
			Args: req.Args,
		}

		if err := s.services.Repository.MutatePublication(req.Id, nil, mut); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("failed to mutate publication %s at line %d: %s", req.Id, seq, err).Error())
			if err = stream.Send(&api.MutateResponse{
				Response: &api.MutateResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to mutate publication: %v", err)
			}
			continue
		}

		if err = stream.Send(&api.MutateResponse{
			Response: &api.MutateResponse_Message{
				Message: fmt.Sprintf("mutated publication %s at line %d", req.Id, seq),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to mutate publication: %v", err)
		}
	}
}

func (s *server) GetPublicationHistory(req *api.GetPublicationHistoryRequest, stream api.Biblio_GetPublicationHistoryServer) (err error) {
	var callbackErr error
	streamErr := s.services.Repository.PublicationHistory(req.Id, func(p *models.Publication) bool {
		j, err := json.Marshal(p)
		if err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.GetPublicationHistoryResponse{
				Response: &api.GetPublicationHistoryResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				callbackErr = err
				return false
			}
		} else {
			if err = stream.Send(&api.GetPublicationHistoryResponse{
				Response: &api.GetPublicationHistoryResponse_Publication{
					Publication: &api.Publication{
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
		return status.Errorf(codes.Internal, "could not get publication history: %v", callbackErr)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "could not get publication history: %v	", streamErr)
	}

	return nil
}

func (s *server) PurgePublication(ctx context.Context, req *api.PurgePublicationRequest) (*api.PurgePublicationResponse, error) {
	_, err := s.services.Repository.GetPublication(req.Id)

	if err != nil {
		if errors.Is(err, backends.ErrNotFound) {
			grpcErr := status.New(codes.NotFound, fmt.Errorf("could not find publication with id %s", req.Id).Error())
			return &api.PurgePublicationResponse{
				Response: &api.PurgePublicationResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %v", err)
		}
	}

	// TODO purgePublication doesn't return an error if the record for req.Id can't be found
	if err := s.services.Repository.PurgePublication(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication with id %s: %s", req.Id, err)
	}

	// TODO this will complain if the above didn't throw a 'not found' error
	if err := s.services.PublicationSearchService.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication from index with id %s: %s", req.Id, err)
	}

	return &api.PurgePublicationResponse{
		Response: &api.PurgePublicationResponse_Ok{
			Ok: true,
		},
	}, nil
}

func (s *server) PurgeAllPublications(ctx context.Context, req *api.PurgeAllPublicationsRequest) (*api.PurgeAllPublicationsResponse, error) {
	if !req.Confirm {
		grpcErr := status.New(codes.Internal, fmt.Errorf("confirm property in request is not set to true").Error())
		return &api.PurgeAllPublicationsResponse{
			Response: &api.PurgeAllPublicationsResponse_Error{
				Error: grpcErr.Proto(),
			},
		}, nil
	}

	if err := s.services.Repository.PurgeAllPublications(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all publications: %s", err)
	}

	if err := s.services.PublicationSearchService.DeleteAll(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete publication from index: %w", err)
	}

	return &api.PurgeAllPublicationsResponse{
		Response: &api.PurgeAllPublicationsResponse_Ok{
			Ok: true,
		},
	}, nil
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

		p := &models.Publication{}
		if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not read json input: %s", err).Error())
			if err = stream.Send(&api.ValidatePublicationsResponse{
				Response: &api.ValidatePublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to validate publications: %v", err)
			}
			continue
		}

		err = p.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			if err = stream.Send(&api.ValidatePublicationsResponse{
				Response: &api.ValidatePublicationsResponse_Results{
					Results: &api.ValidateResults{
						Seq:     seq,
						Id:      p.ID,
						Message: validationErrs.Error(),
					},
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to validate publications: %v", err)
			}
		} else if err != nil {
			return status.Errorf(codes.Internal, "failed to validate publications: %v", err)
		}
	}
}

func (s *server) ReindexPublications(req *api.ReindexPublicationsRequest, stream api.Biblio_ReindexPublicationsServer) error {
	startTime := time.Now()
	indexed := 0

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: "Indexing to a new index",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	var swErr error
	var swIdxErr error

	switcher, err := s.services.PublicationSearchService.NewIndexSwitcher(backends.BulkIndexerConfig{
		OnError: func(err error) {
			grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index publication: %v", err).Error())
			if err = stream.Send(&api.ReindexPublicationsResponse{
				Response: &api.ReindexPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				swErr = err
			}
		},
		OnIndexError: func(id string, err error) {
			grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index publication %s: %w", id, err).Error())
			if err = stream.Send(&api.ReindexPublicationsResponse{
				Response: &api.ReindexPublicationsResponse_Error{
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

	streamErr := s.services.Repository.EachPublication(func(p *models.Publication) bool {
		if err := switcher.Index(ctx, p); err != nil {
			grpcErr := status.New(codes.Internal, fmt.Errorf("indexing failed for publication [id: %s] : %s", p.ID, err).Error())
			if err = stream.Send(&api.ReindexPublicationsResponse{
				Response: &api.ReindexPublicationsResponse_Error{
					Error: grpcErr.Proto(),
				},
			}); err != nil {
				// PICK UP ERROR
				return false
			}
		}

		indexed++

		return true
	})

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	if swErr != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", swErr)
	}

	if swIdxErr != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", swIdxErr)
	}

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: fmt.Sprintf("Indexed %d publications...", indexed),
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: "Switching to new index...",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	if err := switcher.Switch(ctx); err != nil {
		return status.Errorf(codes.Internal, "indexing failed: %v", err)
	}

	endTime := time.Now()

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: "Indexing changes since start of reindex...",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	for {
		indexed = 0

		var biErr error
		var biIdxErr error

		bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
			OnError: func(err error) {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index publication: %v", err).Error())
				if err = stream.Send(&api.ReindexPublicationsResponse{
					Response: &api.ReindexPublicationsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					biErr = err
				}
			},
			OnIndexError: func(id string, err error) {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to index publication %s: %w", id, err).Error())
				if err = stream.Send(&api.ReindexPublicationsResponse{
					Response: &api.ReindexPublicationsResponse_Error{
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
		streamErr := s.services.Repository.PublicationsBetween(startTime, endTime, func(p *models.Publication) bool {
			if err := bi.Index(ctx, p); err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("indexing failed for publication [id: %s] : %s", p.ID, err).Error())
				if err = stream.Send(&api.ReindexPublicationsResponse{
					Response: &api.ReindexPublicationsResponse_Error{
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
			return status.Errorf(codes.Internal, "failed to index publications: %v", callbackErr)
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

		if err := stream.Send(&api.ReindexPublicationsResponse{
			Response: &api.ReindexPublicationsResponse_Message{
				Message: fmt.Sprintf("Indexed %d publications...", indexed),
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to index publications: %v", err)
		}

		startTime = endTime
		endTime = time.Now()
	}

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: "Done",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	return nil
}

func (s *server) TransferPublications(req *api.TransferPublicationsRequest, stream api.Biblio_TransferPublicationsServer) error {
	source := req.Src
	dest := req.Dest

	p, err := s.services.PersonService.GetPerson(dest)
	if err != nil {
		return status.Errorf(codes.Internal, "could not retrieve person %s: %v", dest, err)
	}

	c := &models.Contributor{}
	c.ID = p.ID
	c.FirstName = p.FirstName
	c.LastName = p.LastName
	c.FullName = p.FullName
	c.UGentID = p.UGentID
	c.ORCID = p.ORCID

	for _, pd := range p.Department {
		newDep := models.ContributorDepartment{ID: pd.ID}
		org, orgErr := s.services.OrganizationService.GetOrganization(pd.ID)
		if orgErr == nil {
			newDep.Name = org.Name
		}
		c.Department = append(c.Department, newDep)
	}

	var callbackErr error

	callback := func(p *models.Publication) bool {
		fixed := false

		if p.User != nil {
			if p.User.ID == source {
				p.User = &models.PublicationUser{
					ID:   c.ID,
					Name: c.FullName,
				}

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: user: %s -> %s", p.ID, p.SnapshotID, source, c.ID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		if p.Creator != nil {
			if p.Creator.ID == source {
				p.Creator = &models.PublicationUser{
					ID:   c.ID,
					Name: c.FullName,
				}

				if len(c.Department) > 0 {
					org, orgErr := s.services.OrganizationService.GetOrganization(c.Department[0].ID)
					if orgErr != nil {
						callbackErr = fmt.Errorf("p: %s: s: %s ::: creator: could not fetch department for %s: %v", p.ID, p.SnapshotID, c.ID, orgErr)
						return false
					} else {
						p.AddDepartmentByOrg(org)
					}
				}

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: creator: %s -> %s", p.ID, p.SnapshotID, source, c.ID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, a := range p.Author {
			if a.ID == source {
				p.SetContributor("author", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: author: %s -> %s", p.ID, p.SnapshotID, a.ID, c.ID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, e := range p.Editor {
			if e.ID == source {
				p.SetContributor("editor", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: editor: %s -> %s", p.ID, p.SnapshotID, e.ID, c.ID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, s := range p.Supervisor {
			if s.ID == source {
				p.SetContributor("supervisor", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: supervisor: %s -> %s", p.ID, p.SnapshotID, s.ID, c.ID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		if fixed {
			errUpdate := s.services.Repository.UpdatePublicationInPlace(p)
			if errUpdate != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("p: %s: s: %s ::: could not update snapshot: %s", p.ID, p.SnapshotID, errUpdate).Error())
				if err = stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}
			}
		}

		return true
	}

	var streamErr error

	if req.Publicationid != "" {
		streamErr = s.services.Repository.PublicationHistory(req.Publicationid, callback)
	} else {
		streamErr = s.services.Repository.EachPublicationSnapshot(callback)
	}

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "failed to transfer publication: %v", callbackErr)
	}

	if streamErr != nil {
		return status.Errorf(codes.Internal, "failed to transfer publication: %v", streamErr)
	}

	return nil
}

func (s *server) CleanupPublications(req *api.CleanupPublicationsRequest, stream api.Biblio_CleanupPublicationsServer) error {
	var callbackErr error

	count := 0
	streamErr := s.services.Repository.EachPublication(func(p *models.Publication) bool {
		// Guard
		fixed := false

		// Add the department "tree" property if it is missing.
		for _, dep := range p.Department {
			if dep.Tree == nil {
				depID := dep.ID
				org, orgErr := s.services.OrganizationService.GetOrganization(depID)
				if orgErr == nil {
					p.RemoveDepartment(depID)
					p.AddDepartmentByOrg(org)
					fixed = true
				}
			}
		}

		// Trim keywords, remove empty keywords
		var cleanKeywords []string
		for _, kw := range p.Keyword {
			cleanKw := strings.TrimSpace(kw)
			if cleanKw != kw || cleanKw == "" {
				fixed = true
			}
			if cleanKw != "" {
				cleanKeywords = append(cleanKeywords, cleanKw)
			}
		}
		p.Keyword = cleanKeywords

		// Remove empty links (temporary)
		for _, l := range p.Link {
			if l.URL == "" {
				p.RemoveLink(l.ID)
				fixed = true
			}
		}

		// Remove unused fields
		if p.CleanupUnusedFields() {
			fixed = true
		}

		// Save record if changed
		if fixed {
			p.User = nil

			if err := p.Validate(); err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to validate publication[snapshot_id: %s, id: %s] : %v", p.SnapshotID, p.ID, err).Error())
				if err = stream.Send(&api.CleanupPublicationsResponse{
					Response: &api.CleanupPublicationsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				return true
			}

			err := s.services.Repository.UpdatePublication(p.SnapshotID, p, nil)
			if err != nil {
				grpcErr := status.New(codes.Internal, fmt.Errorf("failed to update publication[snapshot_id: %s, id: %s] : %v", p.SnapshotID, p.ID, err).Error())
				if err = stream.Send(&api.CleanupPublicationsResponse{
					Response: &api.CleanupPublicationsResponse_Error{
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
				grpcErr := status.New(codes.Internal, fmt.Errorf("conflict detected for publication[snapshot_id: %s, id: %s] : %v", p.SnapshotID, p.ID, err).Error())
				if err = stream.Send(&api.CleanupPublicationsResponse{
					Response: &api.CleanupPublicationsResponse_Error{
						Error: grpcErr.Proto(),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				return true
			}

			if err := stream.Send(&api.CleanupPublicationsResponse{
				Response: &api.CleanupPublicationsResponse_Message{
					Message: fmt.Sprintf("fixed publication[snapshot_id: %s, id: %s]", p.SnapshotID, p.ID),
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

	if err := stream.Send(&api.CleanupPublicationsResponse{
		Response: &api.CleanupPublicationsResponse_Message{
			Message: fmt.Sprintf("done. cleaned %d publications.", count),
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to clean up publications: %v", err)
	}

	return nil
}

// SyncPublicationContributors synchronizes data stored in the authority table "people" with the related
// records as stored in publications under "author", "supervisor" and "editor", which we all call "contributors".
// More specifically, every contributor in a publication that has an attribute "id",
// is expected to have a related record with that id in the authority table.
//
// The attributes that are copied from the authority table are:
// - first_name
// - last_name
// (- full_name) (TODO: see below)
// - orcid
// - ugent id's
// - department id's (and their name)
func (s *server) SyncPublicationContributors(req *api.SyncPublicationContributorsRequest, stream api.Biblio_SyncPublicationContributorsServer) error {
	repository := s.services.Repository
	personService := s.services.PersonService
	orgService := s.services.OrganizationService

	sendErr := func(stream api.Biblio_SyncPublicationContributorsServer, e error) error {
		grpcErr := status.New(codes.Internal, e.Error())
		return stream.Send(&api.SyncPublicationContributorsResponse{
			Response: &api.SyncPublicationContributorsResponse_Error{
				Error: grpcErr.Proto(),
			},
		})
	}

	var callbackErr error
	// cache does not store nil values
	invalidContributors := map[string]error{}
	contributorRoles := []string{"author", "editor", "supervisor"}

	// startErr: error that stores why the loop could not even start
	startErr := repository.EachPublication(func(p *models.Publication) bool {

		changes := make([]*api.ContributorChange, 0)

		// prefetch all contributors from personService
		// ignore those contributors that did not return anything in a previous call
		contributorIds := make([]string, 0)
		for _, role := range contributorRoles {
			for _, contributor := range p.Contributors(role) {
				if contributor.ID == "" {
					continue
				}
				if _, ok := invalidContributors[contributor.ID]; ok {
					continue
				}
				found := false
				for _, contributorId := range contributorIds {
					if contributorId == contributor.ID {
						found = true
					}
				}
				if !found {
					contributorIds = append(contributorIds, contributor.ID)
				}
			}
		}

		var persons []*models.Person = make([]*models.Person, 0)
		if len(contributorIds) > 0 {

			var personErr error
			persons, personErr = personService.GetPersons(contributorIds)
			if personErr != nil {

				e := sendErr(stream, fmt.Errorf("unable to prefetch contributors for publication with id '%s': %s", p.ID, personErr))

				if e != nil {
					callbackErr = e
					// unable to send error: stop the loop
					return false
				}

				// TODO: send error
				return false
			}

		}

		for _, role := range contributorRoles {

			contributors := p.Contributors(role)

			for _, c := range contributors {

				// only handle records from authority database
				if c.ID == "" {
					continue
				}

				// contributors already handled
				if cErr, ok := invalidContributors[c.ID]; ok {
					e := sendErr(stream, fmt.Errorf("[duplicate] unable to fetch person with id '%s' with role '%s' for publication with id '%s': %s", c.ID, role, p.ID, cErr))

					if e != nil {
						callbackErr = e
						// unable to send error: stop the loop
						return false
					}

					// this is not fatal error, so keep going
					continue
				}

				var person *models.Person
				for _, p := range persons {
					if p.ID == c.ID {
						person = p
						break
					}
				}
				var personErr error
				if person == nil {
					personErr = backends.ErrNotFound
				}

				// person is gone from the authority table
				if personErr != nil {

					invalidContributors[c.ID] = personErr

					e := sendErr(stream, fmt.Errorf("unable to fetch person with id '%s' with role '%s' for publication with id '%s': %s", c.ID, role, p.ID, personErr))

					if e != nil {
						callbackErr = e
						// unable to send error: stop the loop
						return false
					}

					// this is not fatal error, so keep going
					continue
				}

				if c.FirstName != person.FirstName {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "contributor.first_name",
						From:            c.FirstName,
						To:              person.FirstName,
					})
					c.FirstName = person.FirstName
				}
				if c.LastName != person.LastName {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "contributor.last_name",
						From:            c.LastName,
						To:              person.LastName,
					})
					c.LastName = person.LastName
				}

				if c.FullName != person.FullName {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "contributor.full_name",
						From:            c.FullName,
						To:              person.FullName,
					})
					c.FullName = person.FullName
				}

				if c.ORCID != person.ORCID {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "contributor.orcid",
						From:            c.ORCID,
						To:              person.ORCID,
					})
					c.ORCID = person.ORCID
				}

				if !reflect.DeepEqual(c.UGentID, person.UGentID) {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "contributor.ugent_id",
						From:            strings.Join(c.UGentID, ","),
						To:              strings.Join(person.UGentID, ","),
					})
					c.UGentID = append([]string{}, person.UGentID...)
				}

				oldDeps := make([]string, 0, len(c.Department))
				for _, dep := range c.Department {
					oldDeps = append(oldDeps, dep.ID)
				}
				sort.Strings(oldDeps)

				newDeps := make([]string, 0, len(person.Department))
				for _, dep := range person.Department {
					newDeps = append(newDeps, dep.ID)
				}
				sort.Strings(newDeps)

				if !reflect.DeepEqual(oldDeps, newDeps) {
					changes = append(changes, &api.ContributorChange{
						PublicationId:   p.ID,
						ContributorId:   c.ID,
						ContributorRole: role,
						Attribute:       "department",
						From:            strings.Join(oldDeps, ","),
						To:              strings.Join(newDeps, ","),
					})
					c.Department = make([]models.ContributorDepartment, 0)
					for _, pd := range person.Department {
						newDep := models.ContributorDepartment{ID: pd.ID}
						org, orgErr := orgService.GetOrganization(pd.ID)
						if orgErr == nil {
							newDep.Name = org.Name
						} else {
							err := sendErr(stream, fmt.Errorf("unable to fetch organization with id '%s' for publication with id '%s': %s", pd.ID, p.ID, orgErr))
							if err != nil {
								callbackErr = err
								// unable to send error: stop the loop
								return false
							}
						}
						c.Department = append(c.Department, newDep)
					}
				}
			}
		}

		if !req.Noop {
			for _, change := range changes {
				change.Executed = true
			}
		}

		if err := p.Validate(); err != nil {

			for _, validationErr := range err.(validation.Errors) {
				formattedErr := fmt.Errorf("validation failed for publication %s: %w", p.ID, validationErr)
				if sErr := sendErr(stream, formattedErr); sErr != nil {
					callbackErr = sErr
					// unable to send error: stop the loop
					return false
				}
			}

			// validation failed: continue to next record
			return true

		}

		if !req.Noop && len(changes) > 0 {

			err := s.services.Repository.UpdatePublication(p.SnapshotID, p, nil)

			if err != nil {

				// unable to send error: stop the loop
				e := sendErr(stream, fmt.Errorf("failed to update publication[snapshot_id: %s, id: %s] : %v", p.SnapshotID, p.ID, err))
				if e != nil {
					callbackErr = err
					return false
				}

				// update failed: continue to next record
				return true

			}
		}

		// validation, update and indexation ok: list and send all changes to the client
		for _, change := range changes {
			// unable to send message: stop the loop
			if err := stream.Send(&api.SyncPublicationContributorsResponse{
				Response: &api.SyncPublicationContributorsResponse_ContributorChange{
					ContributorChange: change,
				},
			}); err != nil {
				callbackErr = err
				return false
			}
		}

		// record processed successfully
		return true
	})

	if startErr != nil {
		return status.Errorf(codes.Internal, "unable to start command: %v", startErr)
	}

	if callbackErr != nil {
		return status.Errorf(codes.Internal, "unable to complete command: %s", callbackErr.Error())
	}

	return nil
}
