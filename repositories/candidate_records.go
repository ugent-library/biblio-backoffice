package repositories

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/db"
	"github.com/ugent-library/biblio-backoffice/models"
)

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

func (r *Repo) GetCandidateRecords(ctx context.Context, start int, limit int) (int, []*models.CandidateRecord, error) {
	rows, err := r.queries.GetCandidateRecords(ctx, db.GetCandidateRecordsParams{
		Limit:  int32(limit),
		Offset: int32(start),
	})
	if err != nil {
		return 0, nil, err
	}

	var total int
	recs := make([]*models.CandidateRecord, len(rows))
	for i, row := range rows {
		total = int(row.Total)
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
			StatusPersonID: row.StatusPersonID,
			ImportedID:     row.ImportedID,
		}
		if err := json.Unmarshal(rec.Metadata, rec.Publication); err != nil {
			return 0, nil, err
		}
		for _, fn := range r.config.PublicationLoaders {
			if err := fn(rec.Publication); err != nil {
				return 0, nil, err
			}
		}

		recs[i] = rec
	}
	return total, recs, err
}

func (r *Repo) HasCandidateRecords(ctx context.Context) (bool, error) {
	exists, err := r.queries.HasCandidateRecords(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repo) GetCandidateRecordsByPersonID(ctx context.Context, personID string, start int, limit int, newOnly bool) (int, []*models.CandidateRecord, error) {
	query, _ := json.Marshal([]struct {
		PersonID string `json:"person_id"`
	}{{PersonID: personID}})

	rows, err := r.queries.GetCandidateRecordsByPersonID(ctx, db.GetCandidateRecordsByPersonIDParams{
		Query:   query,
		Limit:   int32(limit),
		Offset:  int32(start),
		NewOnly: newOnly,
	})
	if err != nil {
		return 0, nil, err
	}

	var total int
	recs := make([]*models.CandidateRecord, len(rows))
	for i, row := range rows {
		total = int(row.Total)
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
			StatusPersonID: row.StatusPersonID,
			ImportedID:     row.ImportedID,
		}
		for _, fn := range r.config.CandidateRecordLoaders {
			if err := fn(rec); err != nil {
				return 0, nil, err
			}
		}
		if err := json.Unmarshal(rec.Metadata, rec.Publication); err != nil {
			return 0, nil, err
		}
		for _, fn := range r.config.PublicationLoaders {
			if err := fn(rec.Publication); err != nil {
				return 0, nil, err
			}
		}

		recs[i] = rec
	}
	return total, recs, err
}

func (r *Repo) PersonHasCandidateRecords(ctx context.Context, personID string) (bool, error) {
	query, _ := json.Marshal([]struct {
		PersonID string `json:"person_id"`
	}{{PersonID: personID}})
	exists, err := r.queries.PersonHasCandidateRecords(ctx, query)
	if err != nil {
		return false, err
	}
	return exists, nil
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
		StatusPersonID: row.StatusPersonID,
		ImportedID:     row.ImportedID,
	}

	if err := json.Unmarshal(rec.Metadata, rec.Publication); err != nil {
		return nil, err
	}
	for _, fn := range r.config.PublicationLoaders {
		if err := fn(rec.Publication); err != nil {
			return nil, err
		}
	}

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
		StatusPersonID: row.StatusPersonID,
		ImportedID:     row.ImportedID,
	}

	if err := json.Unmarshal(rec.Metadata, rec.Publication); err != nil {
		return nil, err
	}
	for _, fn := range r.config.PublicationLoaders {
		if err := fn(rec.Publication); err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (r *Repo) RejectCandidateRecord(ctx context.Context, id string) error {
	_, err := r.queries.SetCandidateRecordStatus(ctx, db.SetCandidateRecordStatusParams{Status: "rejected", ID: id})

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
