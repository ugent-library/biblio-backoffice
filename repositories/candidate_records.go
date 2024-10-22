package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/db"
	"github.com/ugent-library/biblio-backoffice/models"
)

type candidateRecordRow struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Metadata       json.RawMessage
	Status         string
	DateCreated    time.Time
	StatusDate     *time.Time
	StatusPersonID *string
	ImportedID     *string
	Total          int
}

func (r *Repo) AddCandidateRecord(ctx context.Context, rec *models.CandidateRecord) error {
	rec.ID = ulid.Make().String()
	params := db.AddCandidateRecordParams{
		ID:             rec.ID,
		SourceName:     rec.SourceName,
		SourceID:       rec.SourceID,
		SourceMetadata: rec.SourceMetadata,
		Type:           rec.Type,
		Metadata:       rec.Metadata,
	}
	_, err := r.queries.AddCandidateRecord(ctx, params)
	return err
}

func (r *Repo) GetCandidateRecords(ctx context.Context, searchArgs *models.SearchArgs) (total int, result []*models.CandidateRecord, err error) {
	candidateRecordRows, err := queryRows[candidateRecordRow](r, ctx, buildQuery(searchArgs))
	if err != nil {
		return 0, nil, err
	}

	total, result, err = r.mapRows(candidateRecordRows, func(row candidateRecordRow) *models.CandidateRecord {
		return &models.CandidateRecord{
			ID:             row.ID,
			SourceName:     row.SourceName,
			SourceID:       row.SourceID,
			Type:           row.Type,
			Metadata:       row.Metadata,
			DateCreated:    row.DateCreated,
			Status:         row.Status,
			Publication:    &models.Publication{},
			StatusDate:     row.StatusDate,
			StatusPersonID: lo.FromPtr(row.StatusPersonID),
			ImportedID:     lo.FromPtr(row.ImportedID),
		}
	})

	return
}

func (r *Repo) HasCandidateRecords(ctx context.Context) (bool, error) {
	exists, err := r.queries.HasCandidateRecords(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repo) PersonHasCandidateRecords(ctx context.Context, personID string) (bool, error) {
	exists, err := r.queries.PersonHasCandidateRecords(ctx, getPersonFilter(personID))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repo) CountPersonCandidateRecords(ctx context.Context, personID string) (int, error) {
	n, err := r.queries.CountPersonCandidateRecords(ctx, getPersonFilter(personID))
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (r *Repo) GetCandidateRecordBySource(ctx context.Context, sourceName string, sourceID string) (*models.CandidateRecord, error) {
	row, err := r.queries.GetCandidateRecordBySource(ctx, db.GetCandidateRecordBySourceParams{
		SourceName: sourceName,
		SourceID:   sourceID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rec := &models.CandidateRecord{
		ID:             row.ID,
		SourceName:     row.SourceName,
		SourceID:       row.SourceID,
		Type:           row.Type,
		Metadata:       row.Metadata,
		DateCreated:    row.DateCreated.Time,
		Status:         row.Status,
		Publication:    &models.Publication{},
		StatusDate:     &row.StatusDate.Time,
		StatusPersonID: lo.FromPtr(row.StatusPersonID),
		ImportedID:     lo.FromPtr(row.ImportedID),
	}

	r.loadCandidateRecord(rec)

	return rec, nil
}

func (r *Repo) GetCandidateRecord(ctx context.Context, id string) (*models.CandidateRecord, error) {
	row, err := r.queries.GetCandidateRecord(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rec := &models.CandidateRecord{
		ID:             row.ID,
		SourceName:     row.SourceName,
		SourceID:       row.SourceID,
		Type:           row.Type,
		Metadata:       row.Metadata,
		DateCreated:    row.DateCreated.Time,
		Status:         row.Status,
		Publication:    &models.Publication{},
		StatusDate:     &row.StatusDate.Time,
		StatusPersonID: lo.FromPtr(row.StatusPersonID),
		ImportedID:     lo.FromPtr(row.ImportedID),
	}

	r.loadCandidateRecord(rec)

	return rec, nil
}

func (r *Repo) RejectCandidateRecord(ctx context.Context, id string, user *models.Person) error {
	_, err := r.queries.SetCandidateRecordStatus(ctx, db.SetCandidateRecordStatusParams{
		Status:         "rejected",
		ID:             id,
		StatusPersonID: &user.ID,
	})

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrNotFound
	case err != nil:
		return err
	}

	return nil
}

func (r *Repo) RestoreCandidateRecord(ctx context.Context, id string, user *models.Person) error {
	_, err := r.queries.SetCandidateRecordStatus(ctx, db.SetCandidateRecordStatusParams{
		Status:         "new",
		ID:             id,
		StatusPersonID: &user.ID,
		ImportedID:     nil,
	})

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrNotFound
	case err != nil:
		return err
	}

	return nil

}

func (r *Repo) ImportCandidateRecordAsPublication(ctx context.Context, id string, user *models.Person) (string, error) {
	rec, err := r.GetCandidateRecord(ctx, id)
	if err != nil {
		return "", err
	}

	rec.Publication.ID = ulid.Make().String()
	rec.Publication.CreatorID = user.ID
	rec.Publication.Creator = user

	err = r.tx(ctx, func(r *Repo) error {
		if err := r.SavePublication(rec.Publication, user); err != nil {
			return err
		}

		if _, err := r.queries.SetCandidateRecordStatus(ctx, db.SetCandidateRecordStatusParams{
			ID:             rec.ID,
			Status:         "imported",
			StatusPersonID: &user.ID,
			ImportedID:     &rec.Publication.ID,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return rec.Publication.ID, nil
}

func (r *Repo) GetCandidateRecordsStatusFacet(ctx context.Context, searchArgs *models.SearchArgs) (models.FacetValues, error) {
	query := getBaseQuery("status AS Value", "COUNT(*) AS Count").
		OrderBy("array_position(ARRAY['new', 'imported', 'rejected'], status)").
		GroupBy("Value")

	query = addQueryFilters(query, searchArgs, "status")

	result, err := queryRows[models.Facet](r, ctx, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) GetCandidateRecordsFacultyFacet(ctx context.Context, searchArgs *models.SearchArgs) (models.FacetValues, error) {
	query := getBaseQuery("jsonb_path_query(metadata, '$.related_organizations.organization_id')->>0 AS Value", "COUNT(*) AS Count").
		OrderBy("Value").
		GroupBy("Value")

	query = addQueryFilters(query, searchArgs, "faculty_id")

	result, err := queryRows[models.Facet](r, ctx, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) GetCandidateRecordsPublicationYearFacet(ctx context.Context, searchArgs *models.SearchArgs) (models.FacetValues, error) {
	query := getBaseQuery("metadata->>'year' AS Value", "COUNT(*) AS Count").
		OrderBy("Value DESC").
		GroupBy("Value")

	query = addQueryFilters(query, searchArgs, "year")

	result, err := queryRows[models.Facet](r, ctx, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getBaseQuery(selectColumns ...string) sq.SelectBuilder {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	if len(selectColumns) == 0 {
		selectColumns = []string{"*"}
	}

	return psql.Select(selectColumns...).From("candidate_records")
}

func buildQuery(searchArgs *models.SearchArgs) sq.SelectBuilder {
	query := getBaseQuery("*", "COUNT(*) OVER() AS total")

	query = addQueryFilters(query, searchArgs)

	sort := "default"
	if len(searchArgs.Sort) > 0 {
		sort = searchArgs.Sort[0]
	}
	switch sort {
	case "added-desc":
		query = query.OrderBy("date_created DESC")

	case "added-asc":
		query = query.OrderBy("date_created ASC")

	case "year-desc":
		query = query.OrderBy("metadata->'year' DESC")

	case "year-asc":
		query = query.OrderBy("metadata->'year' ASC")

	default:
		query = query.OrderBy("array_position(ARRAY['new', 'imported', 'rejected'], status)", "date_created DESC")
	}

	return query.
		Limit(uint64(searchArgs.Limit())).
		Offset(uint64(searchArgs.Offset()))
}

func addQueryFilters(query sq.SelectBuilder, searchArgs *models.SearchArgs, omitFilters ...string) sq.SelectBuilder {
	query = query.
		Where(sq.Or{
			sq.Eq{"status": "new"},
			sq.And{
				sq.Expr("status_date IS NOT NULL"),
				sq.LtOrEq{"EXTRACT(DAY FROM (current_timestamp - status_date))": "90"},
			},
		})

	filters := searchArgs.Filters
	if len(omitFilters) > 0 {
		filters = lo.OmitByKeys(filters, omitFilters)
	}

	for field, filterValue := range filters {
		switch field {
		case "status":
			query = query.Where(sq.Eq{"status": filterValue})

		case "faculty_id":
			conditions := lo.Map(filterValue, func(facultyID string, _ int) sq.Sqlizer {
				return sq.Expr("metadata->'related_organizations' @> ?::jsonb", getFacultyFilter(facultyID))
			})
			query = query.Where(sq.Or(conditions))

		case "year":
			query = query.Where(sq.Eq{"metadata->>'year'": filterValue})

		case "person_id":
			personFilter := getPersonFilter(filterValue[0])
			query = query.Where("(metadata->'author' @> ?::jsonb OR metadata->'supervisor' @> ?::jsonb)", personFilter, personFilter)
		}
	}

	return query
}

func getPersonFilter(personID string) []byte {
	personFilter, _ := json.Marshal([]struct {
		PersonID string `json:"person_id"`
	}{{PersonID: personID}})

	return personFilter
}

func getFacultyFilter(facultyID string) []byte {
	facultyFilter, _ := json.Marshal([]models.RelatedOrganization{{OrganizationID: facultyID}})

	return facultyFilter
}

func (r *Repo) mapRows(rows []candidateRecordRow, iteratee func(candidateRecordRow) *models.CandidateRecord) (int, []*models.CandidateRecord, error) {
	if len(rows) == 0 {
		return 0, make([]*models.CandidateRecord, 0), nil
	}

	results := make([]*models.CandidateRecord, 0, len(rows))

	for _, item := range rows {
		result := iteratee(item)

		err := r.loadCandidateRecord(result)
		if err != nil {
			return 0, nil, err
		}

		results = append(results, result)
	}

	return (rows)[0].Total, results, nil
}

func (r *Repo) loadCandidateRecord(result *models.CandidateRecord) error {
	for _, fn := range r.config.CandidateRecordLoaders {
		if err := fn(result); err != nil {
			return err
		}
	}

	if err := json.Unmarshal(result.Metadata, result.Publication); err != nil {
		return err
	}
	for _, fn := range r.config.PublicationLoaders {
		if err := fn(result.Publication); err != nil {
			return err
		}
	}

	return nil
}
