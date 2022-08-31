package server

import (
	"context"
	"fmt"
	"io"
	"sync"

	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) GetDataset(ctx context.Context, req *api.GetDatasetRequest) (*api.GetDatasetResponse, error) {
	dataset, err := s.services.Repository.GetDataset(req.Id)
	if err != nil {
		// TODO How do we differentiate between errors? e.g. NotFound vs. Internal (database unavailable,...)
		return nil, status.Errorf(codes.Internal, "Could not get dataset with id %d: %w", req.Id, err)
	}

	res := &api.GetDatasetResponse{Dataset: datasetToMessage(dataset)}

	return res, nil
}

func (s *server) GetAllDatasets(req *api.GetAllDatasetsRequest, stream api.Biblio_GetAllDatasetsServer) (err error) {
	return s.services.Repository.EachDataset(func(p *models.Dataset) bool {
		res := &api.GetAllDatasetsResponse{Dataset: datasetToMessage(p)}
		if err = stream.Send(res); err != nil {
			// TODO error handling
			return false
		}
		return true
	})
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
		return nil, status.Errorf(codes.Internal, "Could not search datasets: %s", req.Query, err)
	}

	res := &api.SearchDatasetsResponse{
		Limit:  int32(hits.Limit),
		Offset: int32(hits.Offset),
		Total:  int32(hits.Total),
		Hits:   make([]*api.Dataset, len(hits.Hits)),
	}
	for i, p := range hits.Hits {
		res.Hits[i] = datasetToMessage(p)
	}

	return res, nil
}

func (s *server) UpdataDataset(ctx context.Context, req *api.UpdateDatasetRequest) (*api.UpdateDatasetResponse, error) {
	p := messageToDataset(req.Dataset)

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for dataset %s: %w", p.ID, err)
	}

	if err := s.services.Repository.UpdateDataset(req.Dataset.SnapshotId, p); err != nil {
		// TODO How do we differentiate between errors?
		return nil, status.Errorf(codes.Internal, "failed to store dataset %s, %w", p.ID, err)
	}
	if err := s.services.DatasetSearchService.Index(p); err != nil {
		// TODO How do we differentiate between errors
		return nil, status.Errorf(codes.Internal, "failed to index dataset %s, %w", p.ID, err)
	}

	return &api.UpdateDatasetResponse{}, nil
}

// TODO catch indexing errors
func (s *server) AddDatasets(stream api.Biblio_AddDatasetsServer) error {
	// indexing channel
	indexC := make(chan *models.Dataset)

	var indexWG sync.WaitGroup

	// start bulk indexer
	indexWG.Add(1)
	go func() {
		defer indexWG.Done()
		s.services.DatasetSearchService.IndexMultiple(indexC)
	}()

	defer func() {
		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	}()

	var lineNum int

	for {
		lineNum++

		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read stream: %w", err)
		}

		d := messageToDataset(res.Dataset)

		if d.ID == "" {
			d.ID = ulid.MustGenerate()
		}

		if d.Status == "" {
			d.Status = "private"
		}

		for i, val := range d.Abstract {
			if val.ID == "" {
				val.ID = ulid.MustGenerate()
			}
			d.Abstract[i] = val
		}

		if err := d.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for dataset %s at line %d: %w", d.ID, lineNum, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Messsage: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.SaveDataset(d); err != nil {
			msg := fmt.Errorf("failed to store dataset %s at line %d: %w", d.ID, lineNum, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Messsage: msg}); err != nil {
				return status.Errorf(codes.Internal, msg)
			}
			continue
		}

		indexC <- d
	}
}

func datasetToMessage(d *models.Dataset) *api.Dataset {
	msg := &api.Dataset{}

	msg.Id = d.ID

	switch d.Status {
	case "private":
		msg.Status = api.Dataset_STATUS_PRIVATE
	case "public":
		msg.Status = api.Dataset_STATUS_PUBLIC
	case "deleted":
		msg.Status = api.Dataset_STATUS_DELETED
	case "returned":
		msg.Status = api.Dataset_STATUS_RETURNED
	}

	for _, val := range d.Abstract {
		msg.Abstract = append(msg.Abstract, &api.Text{
			Id:   val.ID,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range d.Author {
		msg.Author = append(msg.Author, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
		})
	}

	// msg.BatchId = d.BatchID

	if d.DateCreated != nil {
		msg.DateCreated = timestamppb.New(*d.DateCreated)
	}
	if d.DateUpdated != nil {
		msg.DateUpdated = timestamppb.New(*d.DateUpdated)
	}
	if d.DateFrom != nil {
		msg.DateFrom = timestamppb.New(*d.DateFrom)
	}
	if d.DateUntil != nil {
		msg.DateUntil = timestamppb.New(*d.DateUntil)
	}

	msg.Title = d.Title

	for _, val := range d.Department {
		msg.Organization = append(msg.Organization, &api.RelatedOrganization{
			Id: val.ID,
		})
	}

	msg.CreatorId = d.CreatorID

	msg.UserId = d.UserID

	msg.Doi = d.DOI

	msg.Keyword = d.Keyword

	msg.Url = d.URL

	msg.Year = d.Year

	msg.ReviewerNote = d.ReviewerNote

	msg.ReviewerTags = d.ReviewerTags

	msg.SnapshotId = d.SnapshotID

	msg.Locked = d.Locked

	msg.Message = d.Message

	msg.AccessLevel = d.AccessLevel

	msg.Format = d.Format

	msg.License = d.License

	// msg.Publication = d.Publication

	msg.Publisher = d.Publisher

	for _, val := range d.Project {
		msg.Project = append(msg.Project, &api.RelatedProject{
			Id: val.ID,
		})
	}

	for _, val := range d.RelatedPublication {
		msg.Publication = append(msg.Publication, &api.RelatedPublication{
			Id: val.ID,
		})
	}

	return msg
}

func messageToDataset(msg *api.Dataset) *models.Dataset {
	d := &models.Dataset{}

	d.ID = msg.Id

	switch msg.Status {
	case api.Dataset_STATUS_PRIVATE:
		d.Status = "private"
	case api.Dataset_STATUS_PUBLIC:
		d.Status = "public"
	case api.Dataset_STATUS_DELETED:
		d.Status = "deleted"
	case api.Dataset_STATUS_RETURNED:
		d.Status = "returned"
	}

	for _, val := range msg.Abstract {
		d.Abstract = append(d.Abstract, models.Text{
			ID:   val.Id,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range msg.Author {
		d.Author = append(d.Author, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
		})
	}

	// d.BatchID = msg.BatchId

	if msg.DateCreated != nil {
		t := msg.DateCreated.AsTime()
		d.DateCreated = &t
	}
	if msg.DateUpdated != nil {
		t := msg.DateUpdated.AsTime()
		d.DateUpdated = &t
	}
	if msg.DateFrom != nil {
		t := msg.DateFrom.AsTime()
		d.DateFrom = &t
	}
	if msg.DateUntil != nil {
		t := msg.DateUntil.AsTime()
		d.DateUntil = &t
	}

	d.Title = msg.Title
	for _, val := range msg.Organization {
		// TODO add tree
		d.Department = append(d.Department, models.DatasetDepartment{
			ID: val.Id,
		})
	}

	d.CreatorID = msg.CreatorId

	d.UserID = msg.UserId

	d.DOI = msg.Doi

	d.Keyword = msg.Keyword

	d.URL = msg.Url

	d.Year = msg.Year

	d.ReviewerNote = msg.ReviewerNote

	d.ReviewerTags = msg.ReviewerTags

	d.SnapshotID = msg.SnapshotId

	d.Locked = msg.Locked

	d.Message = msg.Message

	d.AccessLevel = msg.AccessLevel

	d.Format = msg.Format

	d.License = msg.License

	// d.Publicaiton = msg.Publication

	d.Publisher = msg.Publisher

	for _, val := range msg.Project {
		// TODO add Name
		d.Project = append(d.Project, models.DatasetProject{
			ID: val.Id,
		})
	}

	for _, val := range msg.Publication {
		d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{
			ID: val.Id,
		})
	}

	return d
}
