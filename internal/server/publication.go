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
			return nil, status.Errorf(codes.NotFound, "could not find publication with id %s", req.Id)
		} else {
			return nil, status.Errorf(codes.Internal, "could not get publication with id %s: %v", err)
		}
	}

	j, err := json.Marshal(p)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not unmarshal publication with id %s: %v", err)
	}

	apip := &api.Publication{
		Payload: j,
	}

	res := &api.GetPublicationResponse{
		Publication: apip,
	}

	return res, nil
}

func (s *server) GetAllPublications(req *api.GetAllPublicationsRequest, stream api.Biblio_GetAllPublicationsServer) (err error) {
	// TODO make this a cancelible context which breaks the EachPublication loop when the client goes away
	ctx := context.TODO()

	errorStream := s.services.Repository.EachPublication(ctx, func(p *models.Publication) error {
		j, err := json.Marshal(p)
		if err != nil {
			return err
		}
		apip := &api.Publication{
			Payload: j,
		}
		res := &api.GetAllPublicationsResponse{Publication: apip}
		if err = stream.Send(res); err != nil {
			return err
		}
		return nil
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get all publications: %v", errorStream)
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
	ctx := context.TODO()
	errorStream := s.services.Repository.PublicationHistory(ctx, req.Id, func(p *models.Publication) error {
		j, err := json.Marshal(p)
		if err != nil {
			return err
		}

		apip := &api.Publication{
			Payload: j,
		}

		res := &api.GetPublicationHistoryResponse{Publication: apip}
		if err = stream.Send(res); err != nil {
			return err
		}

		return nil
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get publication history: %v	", errorStream)
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
	startTime := time.Now()
	indexed := 0

	msg := "Indexing to new index..."
	if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
		return err
	}

	var biErr error
	// var biIndexErr error
	switcher, err := s.services.PublicationSearchService.NewIndexSwitcher(backends.BulkIndexerConfig{
		OnError: func(err error) {
			biErr = fmt.Errorf("indexing failed : %s", err)
		},
		OnIndexError: func(id string, err error) {
			// TODO review: does this even work?
			// biIndexErr = fmt.Errorf("indexing failed for publication [id: %s] : %s", id, err)
		},
	})

	if err != nil {
		return status.Errorf(codes.Internal, "indexing failed: %v", err)
	}

	ctx := stream.Context()
	s.services.Repository.EachPublication(ctx, func(p *models.Publication) error {
		if err := switcher.Index(ctx, p); err != nil {
			msg = fmt.Sprintf("indexing failed for publication [id: %s] : %s", p.ID, err)
			if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
				return err
			}
		}
		indexed++
		return nil
	})

	msg = fmt.Sprintf("Indexed %d publications...", indexed)
	if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
		return err
	}

	msg = "Switching to new index..."
	if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
		return err
	}

	if err := switcher.Switch(ctx); err != nil {
		return status.Errorf(codes.Internal, "indexing failed: %v", err)
	}

	if biErr != nil {
		// TODO
	}

	// if biIndexErr != nil {
	// 	// TODO
	// }

	endTime := time.Now()

	msg = "Indexing changes since start of reindex..."
	if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
		return err
	}

	for {
		indexed = 0

		var biErr error
		bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
			OnError: func(err error) {
				biErr = fmt.Errorf("indexing failed : %s", err)
			},
			OnIndexError: func(id string, err error) {
				// TODO: review: does this work properly?
				// errc <- fmt.Errorf("indexing failed for publication [id: %s] : %s", id, err)
			},
		})

		if err != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", err)
		}

		err = s.services.Repository.PublicationsBetween(startTime, endTime, func(p *models.Publication) error {
			if err := bi.Index(ctx, p); err != nil {
				msg = fmt.Sprintf("indexing failed for publication [id: %s] : %s", p.ID, err)
				if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
					return err
				}
			}
			indexed++
			return nil
		})
		if err != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", err)
		}

		if err = bi.Close(ctx); err != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", err)
		}

		if biErr != nil {
			return status.Errorf(codes.Internal, "indexing failed: %v", biErr)
		}

		if indexed == 0 {
			break
		}

		msg = fmt.Sprintf("Indexed %d publications...", indexed)
		if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
			return err
		}

		startTime = endTime
		endTime = time.Now()
	}

	msg = "Done."
	if err := stream.Send(&api.ReindexPublicationsResponse{Message: msg}); err != nil {
		return err
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

	callback := func(p *models.Publication) error {
		fixed := false

		if p.User != nil {
			if p.User.ID == source {
				p.User = &models.PublicationUser{
					ID:   c.ID,
					Name: c.FullName,
				}

				msg := fmt.Sprintf("p: %s: s: %s ::: user: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
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
						return fmt.Errorf("p: %s: s: %s ::: creator: could not fetch department for %s: %v", p.ID, p.SnapshotID, c.ID, orgErr)
					} else {
						p.AddDepartmentByOrg(org)
					}
				}

				msg := fmt.Sprintf("p: %s: s: %s ::: creator: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
				}
				fixed = true
			}
		}

		for k, a := range p.Author {
			if a.ID == source {
				p.SetContributor("author", k, c)
				msg := fmt.Sprintf("p: %s: s: %s ::: author: %s -> %s", p.ID, p.SnapshotID, a.ID, c.ID)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
				}
				fixed = true
			}
		}

		for k, e := range p.Editor {
			if e.ID == source {
				p.SetContributor("editor", k, c)
				msg := fmt.Sprintf("p: %s: s: %s ::: editor: %s -> %s", p.ID, p.SnapshotID, e.ID, c.ID)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
				}
				fixed = true
			}
		}

		for k, s := range p.Supervisor {
			if s.ID == source {
				p.SetContributor("supervisor", k, c)
				msg := fmt.Sprintf("p: %s: s: %s ::: supervisor: %s -> %s", p.ID, p.SnapshotID, s.ID, c.ID)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
				}
				fixed = true
			}
		}

		if fixed {
			errUpdate := s.services.Repository.UpdatePublicationInPlace(p)
			if errUpdate != nil {
				// TODO turn this into a fatal error
				msg := fmt.Sprintf("p: %s: s: %s ::: Could not update snapshot: %s", p.ID, p.SnapshotID, errUpdate)
				if err := stream.Send(&api.TransferPublicationsResponse{Message: msg}); err != nil {
					return err
				}
			}
		}

		return nil
	}

	/* TODO We're altering data, so we want to avoid data races by clients
	   making the same call concurrently. Needs to be revised.
	*/
	// s.mu.transferPublications.Lock()
	// defer s.mu.transferPublications.Unlock()

	ctx := context.TODO()
	if req.Publicationid != "" {
		s.services.Repository.PublicationHistory(ctx, req.Publicationid, callback)
	} else {
		s.services.Repository.EachPublicationSnapshot(ctx, callback)
	}

	return nil
}

func (s *server) CleanupPublications(req *api.CleanupPublicationsRequest, stream api.Biblio_CleanupPublicationsServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bi, err := s.services.PublicationSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
		OnError: func(err error) {
			log.Printf("indexing failed : %s", err)
		},
		OnIndexError: func(id string, err error) {
			log.Printf("indexing failed for publication [id: %s] : %s", id, err)
		},
	})

	if err != nil {
		return status.Errorf(codes.Internal, "bulk indexer failed: %v", err)
	}
	defer bi.Close(ctx)

	count := 0
	err = s.services.Repository.EachPublication(ctx, func(p *models.Publication) error {
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

		// Save record if changed
		if fixed {
			p.User = nil

			if err := p.Validate(); err != nil {
				msg := fmt.Sprintf(
					"Validation failed for publication[snapshot_id: %s, id: %s] : %v",
					p.SnapshotID,
					p.ID,
					err,
				)

				if err := stream.Send(&api.CleanupPublicationsResponse{Message: msg}); err != nil {
					return err
				}

				return nil
			}

			err := s.services.Repository.UpdatePublication(p.SnapshotID, p, nil)
			if err != nil {
				log.Println(err)
			}

			var conflict *snapstore.Conflict
			if errors.As(err, &conflict) {
				msg := fmt.Sprintf(
					"Conflict detected for publication[snapshot_id: %s, id: %s] : %v",
					p.SnapshotID,
					p.ID,
					err,
				)

				if err := stream.Send(&api.CleanupPublicationsResponse{Message: msg}); err != nil {
					return err
				}

				return nil
			}

			msg := fmt.Sprintf(
				"Fixed publication[snapshot_id: %s, id: %s]",
				p.SnapshotID,
				p.ID,
			)

			if err := stream.Send(&api.CleanupPublicationsResponse{Message: msg}); err != nil {
				return err
			}

			if err := bi.Index(ctx, p); err != nil {
				msg := fmt.Sprintf("indexing failed for publication [id: %s] : %s", p.ID, err)
				if err := stream.Send(&api.CleanupPublicationsResponse{Message: msg}); err != nil {
					return err
				}
			}

			count += 1
		}

		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "could not complete cleanup: %v", err)
	}

	msg := fmt.Sprintf("done. cleaned %d publications.", count)
	if err := stream.Send(&api.CleanupPublicationsResponse{Message: msg}); err != nil {
		return err
	}

	return nil
}
