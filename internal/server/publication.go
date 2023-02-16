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
