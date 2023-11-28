package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"slices"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/okay"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var reSplit = regexp.MustCompile(`\s*[,;]\s*`)

func (s *server) GetPublication(ctx context.Context, req *api.GetPublicationRequest) (*api.GetPublicationResponse, error) {
	p, err := s.services.Repo.GetPublication(req.Id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			grpcErr := status.New(codes.InvalidArgument, fmt.Errorf("could not find publication with id %s", req.Id).Error())
			return &api.GetPublicationResponse{
				Response: &api.GetPublicationResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %v", req.Id, err)
		}
	}

	j, err := json.Marshal(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal publication with id %s: %v", req.Id, err)
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

	streamErr := s.services.Repo.EachPublication(func(p *models.Publication) bool {
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
	hits, err := s.services.PublicationSearchIndex.Search(args)
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

	err := s.services.Repo.UpdatePublication(p.SnapshotID, p, nil)

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

		if err := s.services.Repo.SavePublication(p, nil); err != nil {
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

		if err := s.services.Repo.ImportPublication(p); err != nil {
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

		mut := repositories.Mutation{
			Op:   req.Op,
			Args: req.Args,
		}

		if err := s.services.Repo.MutatePublication(req.Id, nil, mut); err != nil {
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
	streamErr := s.services.Repo.PublicationHistory(req.Id, func(p *models.Publication) bool {
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
	_, err := s.services.Repo.GetPublication(req.Id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			grpcErr := status.New(codes.NotFound, fmt.Errorf("could not find publication with id %s", req.Id).Error())
			return &api.PurgePublicationResponse{
				Response: &api.PurgePublicationResponse_Error{
					Error: grpcErr.Proto(),
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %s", req.Id, err)
		}
	}

	// TODO purgePublication doesn't return an error if the record for req.Id can't be found
	if err := s.services.Repo.PurgePublication(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication with id %s: %s", req.Id, err)
	}

	// TODO this will complain if the above didn't throw a 'not found' error
	if err := s.services.PublicationSearchIndex.Delete(req.Id); err != nil {
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

	if err := s.services.Repo.PurgeAllPublications(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all publications: %s", err)
	}

	if err := s.services.PublicationSearchIndex.DeleteAll(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete publication from index: %s", err)
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
		var validationErrs *okay.Errors
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
	reported := 0

	if err := stream.Send(&api.ReindexPublicationsResponse{
		Response: &api.ReindexPublicationsResponse_Message{
			Message: "Indexing to a new index",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	var swErr error
	var swIdxErr error

	switcher, err := s.services.SearchService.NewPublicationIndexSwitcher(backends.BulkIndexerConfig{
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

	streamErr := s.services.Repo.EachPublication(func(p *models.Publication) bool {
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

		// progress message
		if indexed-500 == reported {
			stream.Send(&api.ReindexPublicationsResponse{
				Response: &api.ReindexPublicationsResponse_Message{
					Message: fmt.Sprintf("Indexing %d publications...", indexed),
				},
			})
			reported = indexed
		}

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
			Message: fmt.Sprintf("Indexed %d publications", indexed),
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
			Message: "Indexing changes since start of reindex",
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to index publications: %v", err)
	}

	for {
		indexed = 0

		var biErr error
		var biIdxErr error

		bi, err := s.services.SearchService.NewPublicationBulkIndexer(backends.BulkIndexerConfig{
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
		streamErr := s.services.Repo.PublicationsBetween(startTime, endTime, func(p *models.Publication) bool {
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
				Message: fmt.Sprintf("Indexed %d publications", indexed),
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

	person, err := s.services.PersonService.GetPerson(dest)
	if err != nil {
		return status.Errorf(codes.Internal, "could not retrieve person %s: %v", dest, err)
	}

	c := &models.Contributor{
		PersonID: person.ID,
		Person:   person,
	}

	var callbackErr error

	callback := func(p *models.Publication) bool {
		fixed := false

		if p.UserID != "" {
			if p.UserID == source {
				p.UserID = person.ID
				p.User = person

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: user: %s -> %s", p.ID, p.SnapshotID, source, c.PersonID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		if p.CreatorID != "" {
			if p.CreatorID == source {
				p.CreatorID = person.ID
				p.Creator = person

				if len(c.Person.Affiliations) > 0 {
					org, orgErr := s.services.OrganizationService.GetOrganization(c.Person.Affiliations[0].OrganizationID)
					if orgErr != nil {
						callbackErr = fmt.Errorf("p: %s: s: %s ::: creator: could not fetch department for %s: %v", p.ID, p.SnapshotID, c.PersonID, orgErr)
						return false
					} else {
						p.AddOrganization(org)
					}
				}

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: creator: %s -> %s", p.ID, p.SnapshotID, source, c.PersonID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, a := range p.Author {
			if a.PersonID == source {
				p.SetContributor("author", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: author: %s -> %s", p.ID, p.SnapshotID, a.PersonID, c.PersonID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, e := range p.Editor {
			if e.PersonID == source {
				p.SetContributor("editor", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: editor: %s -> %s", p.ID, p.SnapshotID, e.PersonID, c.PersonID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		for k, s := range p.Supervisor {
			if s.PersonID == source {
				p.SetContributor("supervisor", k, c)

				if err := stream.Send(&api.TransferPublicationsResponse{
					Response: &api.TransferPublicationsResponse_Message{
						Message: fmt.Sprintf("p: %s: s: %s ::: supervisor: %s -> %s", p.ID, p.SnapshotID, s.PersonID, c.PersonID),
					},
				}); err != nil {
					callbackErr = err
					return false
				}

				fixed = true
			}
		}

		if fixed {
			errUpdate := s.services.Repo.UpdatePublicationInPlace(p)
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
		streamErr = s.services.Repo.PublicationHistory(req.Publicationid, callback)
	} else {
		streamErr = s.services.Repo.EachPublicationSnapshot(callback)
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
	streamErr := s.services.Repo.EachPublication(func(p *models.Publication) bool {
		// guard
		fixed := false

		// correctly set HasBeenPublic (only needs to run once)
		if p.Status == "deleted" && !p.HasBeenPublic {
			s.services.Repo.PublicationHistory(p.ID, func(pp *models.Publication) bool {
				if pp.Status == "public" {
					p.HasBeenPublic = true
					fixed = true
					return false
				}
				return true
			})
		}

		// remove empty links (only needs to run once)
		for _, l := range p.Link {
			if l.URL == "" {
				p.RemoveLink(l.ID)
				fixed = true
			}
		}

		// set subtype from wos type (only needs to run once)
		if p.WOSType != "" {
			wosTypes := reSplit.Split(p.WOSType, -1)
			for i, t := range wosTypes {
				wosTypes[i] = strings.ToLower(t)
			}
			firstWOSType := wosTypes[0]
			if p.Type == "journal_article" && p.JournalArticleType == "" {
				switch {
				case slices.Contains(wosTypes, "article") && slices.Contains(wosTypes, "proceedings paper"):
					p.JournalArticleType = "proceedingsPaper"
					fixed = true
				case firstWOSType == "journal article" || firstWOSType == "article" || firstWOSType == "journal paper":
					p.JournalArticleType = "original"
					fixed = true
				case firstWOSType == "review":
					p.JournalArticleType = "review"
					fixed = true
				case firstWOSType == "letter" || firstWOSType == "note" || firstWOSType == "letter/note":
					p.JournalArticleType = "letterNote"
					fixed = true
				}
			}
			if p.Type == "conference" && p.ConferenceType == "" {
				switch {
				case firstWOSType == "meeting abstract":
					p.ConferenceType = "abstract"
					fixed = true
				case firstWOSType == "conference proceeding" || firstWOSType == "proceedings paper" || firstWOSType == "conference paper":
					p.ConferenceType = "proceedingsPaper"
					fixed = true
				case firstWOSType == "poster":
					p.ConferenceType = "poster"
					fixed = true
				}
			}
			if p.Type == "miscellaneous" && p.MiscellaneousType == "" {
				switch {
				case firstWOSType == "art exhibit review":
					p.MiscellaneousType = "exhibitionReview"
					fixed = true
				case firstWOSType == "book review":
					p.MiscellaneousType = "bookReview"
					fixed = true
				case firstWOSType == "dance performance review" || firstWOSType == "theatre review" || firstWOSType == "theater review":
					p.MiscellaneousType = "theatreReview"
					fixed = true
				case firstWOSType == "database review" || firstWOSType == "hardware review" || firstWOSType == "software review":
					p.MiscellaneousType = "productReview"
					fixed = true
				case firstWOSType == "editorial material" || firstWOSType == "editorial":
					p.MiscellaneousType = "editorialMaterial"
					fixed = true
				case firstWOSType == "fiction" || firstWOSType == "creative prose" || firstWOSType == "poetry" || firstWOSType == "script":
					p.MiscellaneousType = "artisticWork"
					fixed = true
				case firstWOSType == "film review" || firstWOSType == "tv review" || firstWOSType == "radio review" || firstWOSType == "video review":
					p.MiscellaneousType = "filmReview"
					fixed = true
				case firstWOSType == "music score review" || firstWOSType == "music performance review" || firstWOSType == "record review":
					p.MiscellaneousType = "musicReview"
					fixed = true
				case firstWOSType == "music score":
					p.MiscellaneousType = "musicEdition"
					fixed = true
				case firstWOSType == "news item":
					p.MiscellaneousType = "newsArticle"
					fixed = true
				case firstWOSType == "correction":
					p.MiscellaneousType = "correction"
					fixed = true
				case firstWOSType == "biographical-item" || firstWOSType == "biographical item" || firstWOSType == "item about an individual":
					p.MiscellaneousType = "biography"
					fixed = true
				case firstWOSType == "bibliography":
					p.MiscellaneousType = "bibliography"
					fixed = true
				case firstWOSType == "preprint":
					p.MiscellaneousType = "preprint"
					fixed = true
				case firstWOSType == "other" || firstWOSType == "discussion" || firstWOSType == "slide":
					p.MiscellaneousType = "other"
					fixed = true
				}
			}
		}

		// remove unused fields
		if p.CleanupUnusedFields() {
			fixed = true
		}

		// remove empty strings from string array
		vacuumArray := func(oldVals []string) []string {
			var newVals []string
			for _, val := range oldVals {
				newVal := strings.TrimSpace(val)
				if newVal != "" {
					newVals = append(newVals, newVal)
				}
				if val != newVal || newVal == "" {
					fixed = true
				}
			}
			return newVals
		}

		p.ISBN = vacuumArray(p.ISBN)
		p.EISBN = vacuumArray(p.EISBN)
		p.ISSN = vacuumArray(p.ISSN)
		p.EISSN = vacuumArray(p.EISSN)
		p.AlternativeTitle = vacuumArray(p.AlternativeTitle)
		p.Keyword = vacuumArray(p.Keyword)
		p.Language = vacuumArray(p.Language)
		p.ResearchField = vacuumArray(p.ResearchField)
		p.ReviewerTags = vacuumArray(p.ReviewerTags)
		p.VABBYear = vacuumArray(p.VABBYear)
		for _, author := range p.Author {
			author.CreditRole = vacuumArray(author.CreditRole)
		}
		for _, editor := range p.Editor {
			editor.CreditRole = vacuumArray(editor.CreditRole)
		}
		for _, supervisor := range p.Supervisor {
			supervisor.CreditRole = vacuumArray(supervisor.CreditRole)
		}

		// save record if changed
		if fixed {
			p.UserID = ""
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

			err := s.services.Repo.UpdatePublication(p.SnapshotID, p, nil)
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
