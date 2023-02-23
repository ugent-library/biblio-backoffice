package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetPublication(ctx context.Context, req *api.GetPublicationRequest) (*api.GetPublicationResponse, error) {
	p, err := s.services.Repository.GetPublication(req.Id)
	if err != nil {
		// TODO How do we differentiate between errors? e.g. NotFound vs. Internal (database unavailable,...)
		return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %w", req.Id, err)
	}

	j, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	apip := &api.Publication{
		Payload: j,
	}

	res := &api.GetPublicationResponse{Publication: apip}

	return res, nil
}

func (s *server) GetAllPublications(req *api.GetAllPublicationsRequest, stream api.Biblio_GetAllPublicationsServer) (err error) {
	// TODO errors in EachPublication aren't caught and pushed upstream. Returning 'false' in the callback
	//   breaks the loop, but EachPublication will return 'nil'.
	//
	//   Logging during streaming doesn't work / isn't possible. The grpc_zap interceptor is only called when
	// 	 GetAllPublication returns an error.
	errorStream := s.services.Repository.EachPublication(func(p *models.Publication) bool {
		j, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}
		apip := &api.Publication{
			Payload: j,
		}
		res := &api.GetAllPublicationsResponse{Publication: apip}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get all publications: %s", errorStream)
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
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "Could not search publications: %s :: %s", req.Query, err)
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
		log.Fatal(err)
	}

	// TODO Fetch user information via better authentication (no basic auth)
	user := &models.User{
		ID:       "n/a",
		FullName: "system user",
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for publication %s: %s", p.ID, err)
	}

	if err := s.services.Repository.UpdatePublication(p.SnapshotID, p, user); err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "failed to store publication %s, %s", p.ID, err)
	}
	if err := s.services.PublicationSearchService.Index(p); err != nil {
		// TODO How do we differentiate between errors
		return nil, status.Errorf(codes.Internal, "failed to index publication %s, %s", p.ID, err)
	}

	return &api.UpdatePublicationResponse{}, nil
}

func (s *server) AddPublications(stream api.Biblio_AddPublicationsServer) error {
	ctx := context.Background()

	bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			msg := fmt.Errorf("failed to index: %w", err).Error()
			// TODO catch error
			stream.Send(&api.AddPublicationsResponse{Message: msg})
		},
		OnIndexError: func(id string, err error) {
			msg := fmt.Errorf("failed to index publication %s: %w", id, err).Error()
			// TODO catch error
			stream.Send(&api.AddPublicationsResponse{Message: msg})
		},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to start an indexer: %s", err)
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

		p := &models.Publication{}
		if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
			return status.Errorf(codes.InvalidArgument, "could not read json input: %s", err)
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

		// TODO this should return structured messages (see validate)
		if err := p.Validate(); err != nil {
			msg := fmt.Errorf("failed to validate publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.SavePublication(p, nil); err != nil {
			msg := fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := bi.Index(ctx, p); err != nil {
			msg := fmt.Errorf("failed to index publication %s at line %d: %w", p.ID, seq, err).Error()
			if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		msg := fmt.Sprintf("stored and indexed publication %s at line %d", p.ID, seq)
		if err = stream.Send(&api.AddPublicationsResponse{Message: msg}); err != nil {
			return err
		}
	}
}

func (s *server) ImportPublications(stream api.Biblio_ImportPublicationsServer) error {
	ctx := context.Background()

	bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			msg := fmt.Errorf("failed to index: %w", err).Error()
			// TODO catch error
			stream.Send(&api.ImportPublicationsResponse{Message: msg})
		},
		OnIndexError: func(id string, err error) {
			msg := fmt.Errorf("failed to index publication %s: %w", id, err).Error()
			// TODO catch error
			stream.Send(&api.ImportPublicationsResponse{Message: msg})
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

		p := &models.Publication{}
		if err := json.Unmarshal(req.Publication.Payload, p); err != nil {
			log.Fatal(err)
		}

		// TODO this should return structured messages (see validate)
		if err := p.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.ImportPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
		if err := s.services.Repository.ImportCurrentPublication(p); err != nil {
			msg := fmt.Errorf("failed to store publication %s at line %d: %s", p.ID, seq, err).Error()
			if err = stream.Send(&api.ImportPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := bi.Index(ctx, p); err != nil {
			msg := fmt.Errorf("failed to index publication %s at line %d: %w", p.ID, seq, err).Error()
			if err = stream.Send(&api.ImportPublicationsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}
	}
}

func (s *server) GetPublicationHistory(req *api.GetPublicationHistoryRequest, stream api.Biblio_GetPublicationHistoryServer) (err error) {
	errorStream := s.services.Repository.PublicationHistory(req.Id, func(p *models.Publication) bool {
		j, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}
		apip := &api.Publication{
			Payload: j,
		}

		res := &api.GetPublicationHistoryResponse{Publication: apip}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get publication history: %s", errorStream)
	}

	return nil
}

func (s *server) PurgePublication(ctx context.Context, req *api.PurgePublicationRequest) (*api.PurgePublicationResponse, error) {
	if err := s.services.Repository.PurgePublication(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication with id %s: %s", req.Id, err)
	}
	if err := s.services.PublicationSearchService.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge publication from index with id %s: %s", req.Id, err)
	}

	return &api.PurgePublicationResponse{}, nil
}

func (s *server) PurgeAllPublications(ctx context.Context, req *api.PurgeAllPublicationsRequest) (*api.PurgeAllPublicationsResponse, error) {
	if err := s.services.Repository.PurgeAllPublications(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all publications: %s", err)
	}
	if err := s.services.PublicationSearchService.DeleteAll(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete publication index: %w", err)
	}

	return &api.PurgeAllPublicationsResponse{}, nil
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
			log.Fatal(err)
		}

		err = p.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			res := &api.ValidatePublicationsResponse{Seq: seq, Id: p.ID, Message: validationErrs.Error()}
			if err = stream.Send(res); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}

func (s *server) ReindexPublications(req *api.ReindexPublicationsRequest, stream api.Biblio_ReindexPublicationsServer) error {
	msgc := make(chan string, 1)
	errc := make(chan error)

	// cancel() is used to shutdown the async bulk indexer as well
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		startTime := time.Now()
		indexed := 0

		msgc <- "Indexing to new index..."

		switcher, err := s.services.PublicationSearchService.NewIndexSwitcher(backends.BulkIndexerConfig{
			OnError: func(err error) {
				errc <- fmt.Errorf("indexing failed : %s", err)
			},
			OnIndexError: func(id string, err error) {
				errc <- fmt.Errorf("indexing failed for publication [id: %s] : %s", id, err)
			},
		})
		if err != nil {
			errc <- err
		}
		s.services.Repository.EachPublication(func(p *models.Publication) bool {
			if err := switcher.Index(ctx, p); err != nil {
				errc <- fmt.Errorf("indexing failed for publication [id: %s] : %s", p.ID, err)
			}
			indexed++
			return true
		})

		msgc <- fmt.Sprintf("Indexed %d publications...", indexed)

		msgc <- "Switching to new index..."

		if err := switcher.Switch(ctx); err != nil {
			errc <- err
		}

		endTime := time.Now()

		msgc <- "Indexing changes since start of reindex..."

		for {
			indexed = 0

			bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
				OnError: func(err error) {
					errc <- fmt.Errorf("indexing failed : %s", err)
				},
				OnIndexError: func(id string, err error) {
					errc <- fmt.Errorf("indexing failed for publication [id: %s] : %s", id, err)
				},
			})
			if err != nil {
				errc <- err
			}

			err = s.services.Repository.PublicationsBetween(startTime, endTime, func(p *models.Publication) bool {
				if err := bi.Index(ctx, p); err != nil {
					errc <- fmt.Errorf("indexing failed for publication [id: %s] : %s", p.ID, err)
				}
				indexed++
				return true
			})
			if err != nil {
				errc <- err
			}

			if err = bi.Close(ctx); err != nil {
				errc <- err
			}

			if indexed == 0 {
				break
			}

			msgc <- fmt.Sprintf("Indexed %d publications...", indexed)

			startTime = endTime
			endTime = time.Now()
		}

		msgc <- "Done."

		close(msgc)
		close(errc)
	}(ctx)

readChannel:
	for {
		select {
		case err := <-errc:
			return err
		case msg, ok := <-msgc:
			if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
				return err
			}

			if !ok {
				// msgc channel closed, processing done.
				break readChannel
			}
		case <-stream.Context().Done():
			// TODO: better error handling / logging server side
			// The client closed the stream on their end, log as an error
			// deferred cancel() is executed, ensures async bulk indexing stops as well.
			return fmt.Errorf("client closed")
		}
	}

	return nil
}

func (s *server) TransferPublications(req *api.TransferPublicationsRequest, stream api.Biblio_TransferPublicationsServer) error {

	repository := *s.repository

	// TODO: this is going to eventually set a huge number of listeners server-side
	// This needs a clean-up function to avoid memory issues.
	repository.AddPublicationListener(func(p *models.Publication) {
		if p.DateUntil == nil {
			if err := s.services.PublicationSearchService.Index(p); err != nil {
				// TODO: listeners need to bubble up errors so we can send them via gRPC to the caller
				log.Fatalf("error indexing publication %s: %v", p.ID, err)
			}
		}
	})

	msgc := make(chan string, 1)
	errc := make(chan error)

	// cancel() is used to shutdown the async bulk indexer as well
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		source := req.Src
		dest := req.Dest

		p, err := s.services.PersonService.GetPerson(dest)
		if err != nil {
			errc <- fmt.Errorf("fatal: could not retrieve person %s: %s", dest, err)
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

		callback := func(p *models.Publication) bool {
			fixed := false

			if p.User != nil {
				if p.User.ID == source {
					p.User = &models.PublicationUser{
						ID:   c.ID,
						Name: c.FullName,
					}

					msgc <- fmt.Sprintf("p: %s: s: %s ::: user: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
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
							errc <- fmt.Errorf("p: %s: s: %s ::: creator: could not fetch department for %s: %s", p.ID, p.SnapshotID, c.ID, orgErr)
						} else {
							p.AddDepartmentByOrg(org)
						}
					}

					msgc <- fmt.Sprintf("p: %s: s: %s ::: creator: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
					fixed = true
				}
			}

			for k, a := range p.Author {
				if a.ID == source {
					p.SetContributor("author", k, c)
					msgc <- fmt.Sprintf("p: %s: s: %s ::: author: %s -> %s", p.ID, p.SnapshotID, a.ID, c.ID)
					fixed = true
				}
			}

			for k, e := range p.Editor {
				if e.ID == source {
					p.SetContributor("editor", k, c)
					msgc <- fmt.Sprintf("p: %s: s: %s ::: editor: %s -> %s", p.ID, p.SnapshotID, e.ID, c.ID)
					fixed = true
				}
			}

			for k, s := range p.Supervisor {
				if s.ID == source {
					p.SetContributor("supervisor", k, c)
					msgc <- fmt.Sprintf("p: %s: s: %s ::: supervisor: %s -> %s", p.ID, p.SnapshotID, s.ID, c.ID)
					fixed = true
				}
			}

			if fixed {
				errUpdate := repository.UpdatePublicationInPlace(p)
				if errUpdate != nil {
					msgc <- fmt.Sprintf("p: %s: s: %s ::: Could not update snapshot: %s", p.ID, p.SnapshotID, errUpdate)
				}
			}

			return true
		}

		if req.Publicationid != "" {
			// TODO Needs a context so it can be cancelled
			repository.PublicationHistory(req.Publicationid, callback)
		} else {
			// TODO Needs a context so it can be cancelled
			repository.EachPublicationSnapshot(callback)
		}

		close(errc)
		close(msgc)
	}(ctx)

readChannel:
	for {
		select {
		// case err := <-errc:
		// 	return err
		case msg, ok := <-msgc:
			if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
				return err
			}

			if !ok {
				// msgc channel closed, processing done.
				break readChannel
			}
		case <-stream.Context().Done():
			// TODO: better error handling / logging server side
			// The client closed the stream on their end, log as an error
			// deferred cancel() is executed, ensures async bulk indexing stops as well.
			return fmt.Errorf("client closed")
		}
	}

	return nil
}
