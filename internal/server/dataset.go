package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/oklog/ulid/v2"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/validation"
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

	res := &api.GetDatasetResponse{Dataset: DatasetToMessage(dataset)}

	return res, nil
}

func (s *server) GetAllDatasets(req *api.GetAllDatasetsRequest, stream api.Biblio_GetAllDatasetsServer) (err error) {
	// TODO errors in EachDataset aren't caught and pushed upstream. Returning 'false' in the callback
	//   breaks the loop, but EachDataset will return 'nil'.
	ErrorStream := s.services.Repository.EachDataset(func(p *models.Dataset) bool {
		res := &api.GetAllDatasetsResponse{Dataset: DatasetToMessage(p)}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if ErrorStream != nil {
		return status.Errorf(codes.Internal, "could not get all datasets: %w", ErrorStream)
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
		return nil, status.Errorf(codes.Internal, "Could not search datasets: %s", req.Query, err)
	}

	res := &api.SearchDatasetsResponse{
		Limit:  int32(hits.Limit),
		Offset: int32(hits.Offset),
		Total:  int32(hits.Total),
		Hits:   make([]*api.Dataset, len(hits.Hits)),
	}
	for i, p := range hits.Hits {
		res.Hits[i] = DatasetToMessage(p)
	}

	return res, nil
}

func (s *server) UpdateDataset(ctx context.Context, req *api.UpdateDatasetRequest) (*api.UpdateDatasetResponse, error) {
	p := MessageToDataset(req.Dataset)

	// TODO Fetch user information via better authentication (no basic auth)
	user := &models.User{
		ID:       "n/a",
		FullName: "system user",
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for dataset %s: %w", p.ID, err)
	}

	if err := s.services.Repository.UpdateDataset(req.Dataset.SnapshotId, p, user); err != nil {
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
			return status.Errorf(codes.Internal, "failed to read stream: %w", err)
		}

		d := MessageToDataset(req.Dataset)

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
			msg := fmt.Errorf("validation failed for dataset %s at line %d: %w", d.ID, seq, err).Error()
			if err = stream.Send(&api.AddDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.SaveDataset(d, nil); err != nil {
			msg := fmt.Errorf("failed to store dataset %s at line %d: %w", d.ID, seq, err).Error()
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
			return status.Errorf(codes.Internal, "failed to read stream: %w", err)
		}

		d := MessageToDataset(req.Dataset)

		// TODO this should return structured messages (see validate)
		if err := d.Validate(); err != nil {
			msg := fmt.Errorf("validation failed for dataset %s at line %d: %w", d.ID, seq, err).Error()
			if err = stream.Send(&api.ImportDatasetsResponse{Message: msg}); err != nil {
				return err
			}
			continue
		}

		if err := s.services.Repository.ImportCurrentDataset(d); err != nil {
			msg := fmt.Errorf("failed to store dataset %s at line %d: %w", d.ID, seq, err).Error()
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
	errorStream := s.services.Repository.DatasetHistory(req.Id, func(p *models.Dataset) bool {
		res := &api.GetDatasetHistoryResponse{Dataset: DatasetToMessage(p)}
		if err = stream.Send(res); err != nil {
			return false
		}
		return true
	})

	if errorStream != nil {
		return status.Errorf(codes.Internal, "could not get dataset history: %w", errorStream)
	}

	return nil
}

func (s *server) PurgeDataset(ctx context.Context, req *api.PurgeDatasetRequest) (*api.PurgeDatasetResponse, error) {
	if err := s.services.Repository.PurgeDataset(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset with id %d: %w", req.Id, err)
	}
	if err := s.services.DatasetSearchService.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge dataset from index with id %d: %w", req.Id, err)
	}

	return &api.PurgeDatasetResponse{}, nil
}

func (s *server) PurgeAllDatasets(ctx context.Context, req *api.PurgeAllDatasetsRequest) (*api.PurgeAllDatasetsResponse, error) {
	if err := s.services.Repository.PurgeAllDatasets(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not purge all datasets: %w", err)
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
			return status.Errorf(codes.Internal, "failed to read stream: %w", err)
		}

		p := MessageToDataset(req.Dataset)

		err = p.Validate()
		var validationErrs validation.Errors
		if errors.As(err, &validationErrs) {
			res := &api.ValidateDatasetsResponse{Seq: seq, Id: p.ID, Message: validationErrs.Error()}
			if err = stream.Send(res); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}

func DatasetToMessage(d *models.Dataset) *api.Dataset {
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

	msg.HasBeenPublic = d.HasBeenPublic

	for _, val := range d.Abstract {
		msg.Abstract = append(msg.Abstract, &api.Text{
			Id:   val.ID,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range d.Contributor {
		var depts []*api.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, &api.ContributorDepartment{
				Id:   dept.ID,
				Name: dept.Name,
			})
		}
		msg.Contributor = append(msg.Contributor, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	for _, val := range d.Author {
		var depts []*api.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, &api.ContributorDepartment{
				Id:   dept.ID,
				Name: dept.Name,
			})
		}
		msg.Author = append(msg.Author, &api.Contributor{
			Id:         val.ID,
			Orcid:      val.ORCID,
			LocalId:    val.UGentID,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	msg.BatchId = d.BatchID

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
		var depts []*api.DepartmentRef
		for _, dept := range val.Tree {
			depts = append(depts, &api.DepartmentRef{
				Id: dept.ID,
			})
		}
		msg.Department = append(msg.Department, &api.Department{
			Id:   val.ID,
			Tree: depts,
		})
	}

	if d.Creator != nil {
		msg.Creator = &api.RelatedUser{Id: d.Creator.ID, Name: d.Creator.Name}
	}

	if d.User != nil {
		msg.User = &api.RelatedUser{Id: d.User.ID, Name: d.User.Name}
	}

	if d.LastUser != nil {
		msg.LastUser = &api.RelatedUser{Id: d.LastUser.ID, Name: d.LastUser.Name}
	}

	msg.Doi = d.DOI

	msg.Keyword = d.Keyword

	msg.Url = d.URL

	msg.Year = d.Year

	msg.ReviewerNote = d.ReviewerNote

	msg.ReviewerTags = d.ReviewerTags

	msg.SnapshotId = d.SnapshotID

	msg.Locked = d.Locked

	msg.Message = d.Message

	switch d.AccessLevel {
	case "info:eu-repo/semantics/openAccess":
		msg.AccessLevel = api.Dataset_ACCESS_LEVEL_OPEN_ACCESS
	case "info:eu-repo/semantics/embargoedAccess":
		msg.AccessLevel = api.Dataset_ACCESS_LEVEL_EMBARGOED_ACCESS
	case "info:eu-repo/semantics/restrictedAccess":
		msg.AccessLevel = api.Dataset_ACCESS_LEVEL_RESTRICTED_ACCESS
	case "info:eu-repo/semantics/closedAccess":
		msg.AccessLevel = api.Dataset_ACCESS_LEVEL_CLOSED_ACCESS
	}

	switch d.AccessLevelAfterEmbargo {
	case "info:eu-repo/semantics/openAccess":
		msg.AccessLevelAfterEmbargo = api.Dataset_ACCESS_LEVEL_OPEN_ACCESS
	case "info:eu-repo/semantics/restrictedAccess":
		msg.AccessLevelAfterEmbargo = api.Dataset_ACCESS_LEVEL_RESTRICTED_ACCESS
	}

	msg.EmbargoDate = d.EmbargoDate

	msg.Format = d.Format

	msg.License = d.License

	msg.OtherLicense = d.OtherLicense

	// msg.Publication = d.Publication

	msg.Publisher = d.Publisher

	for _, val := range d.Project {
		msg.Project = append(msg.Project, &api.RelatedProject{
			Id:   val.ID,
			Name: val.Name,
		})
	}

	for _, val := range d.RelatedPublication {
		msg.RelatedPublication = append(msg.RelatedPublication, &api.RelatedPublication{
			Id: val.ID,
		})
	}

	return msg
}

func MessageToDataset(msg *api.Dataset) *models.Dataset {
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

	d.HasBeenPublic = msg.HasBeenPublic

	for _, val := range msg.Abstract {
		d.Abstract = append(d.Abstract, models.Text{
			ID:   val.Id,
			Text: val.Text,
			Lang: val.Lang,
		})
	}

	for _, val := range msg.Contributor {
		var depts []models.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, models.ContributorDepartment{
				ID:   dept.Id,
				Name: dept.Name,
			})
		}
		d.Contributor = append(d.Contributor, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	for _, val := range msg.Author {
		var depts []models.ContributorDepartment
		for _, dept := range val.Department {
			depts = append(depts, models.ContributorDepartment{
				ID:   dept.Id,
				Name: dept.Name,
			})
		}
		d.Author = append(d.Author, &models.Contributor{
			ID:         val.Id,
			ORCID:      val.Orcid,
			UGentID:    val.LocalId,
			CreditRole: val.CreditRole,
			FirstName:  val.FirstName,
			LastName:   val.LastName,
			FullName:   val.FullName,
			Department: depts,
		})
	}

	d.BatchID = msg.BatchId

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
	for _, val := range msg.Department {
		var depts []models.DatasetDepartmentRef
		for _, dept := range val.Tree {
			depts = append(depts, models.DatasetDepartmentRef{
				ID: dept.Id,
			})
		}
		d.Department = append(d.Department, models.DatasetDepartment{
			ID:   val.Id,
			Tree: depts,
		})
	}

	if msg.Creator != nil {
		d.Creator = &models.DatasetUser{ID: msg.Creator.Id, Name: msg.Creator.Name}
	}

	if msg.User != nil {
		d.User = &models.DatasetUser{ID: msg.User.Id, Name: msg.User.Name}
	}

	if msg.LastUser != nil {
		d.LastUser = &models.DatasetUser{ID: msg.LastUser.Id, Name: msg.LastUser.Name}
	}

	d.DOI = msg.Doi

	d.Keyword = msg.Keyword

	d.URL = msg.Url

	d.Year = msg.Year

	d.ReviewerNote = msg.ReviewerNote

	d.ReviewerTags = msg.ReviewerTags

	d.SnapshotID = msg.SnapshotId

	d.Locked = msg.Locked

	d.Message = msg.Message

	switch msg.AccessLevel {
	case api.Dataset_ACCESS_LEVEL_OPEN_ACCESS:
		d.AccessLevel = "info:eu-repo/semantics/openAccess"
	case api.Dataset_ACCESS_LEVEL_EMBARGOED_ACCESS:
		d.AccessLevel = "info:eu-repo/semantics/embargoedAccess"
	case api.Dataset_ACCESS_LEVEL_RESTRICTED_ACCESS:
		d.AccessLevel = "info:eu-repo/semantics/restrictedAccess"
	case api.Dataset_ACCESS_LEVEL_CLOSED_ACCESS:
		d.AccessLevel = "info:eu-repo/semantics/closedAccess"
	}

	switch msg.AccessLevelAfterEmbargo {
	case api.Dataset_ACCESS_LEVEL_OPEN_ACCESS:
		d.AccessLevelAfterEmbargo = "info:eu-repo/semantics/openAccess"
	case api.Dataset_ACCESS_LEVEL_RESTRICTED_ACCESS:
		d.AccessLevelAfterEmbargo = "info:eu-repo/semantics/restrictedAccess"
	}

	d.EmbargoDate = msg.EmbargoDate

	d.Format = msg.Format

	d.License = msg.License

	d.OtherLicense = msg.OtherLicense

	d.Publisher = msg.Publisher

	for _, val := range msg.Project {
		// TODO add Name
		d.Project = append(d.Project, models.DatasetProject{
			ID:   val.Id,
			Name: val.Name,
		})
	}

	for _, val := range msg.RelatedPublication {
		d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{
			ID: val.Id,
		})
	}

	return d
}
