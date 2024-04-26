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

func (r *Repo) GetCandidateRecords(ctx context.Context, start int, limit int) ([]*models.CandidateRecord, error) {
	rows, err := r.queries.GetCandidateRecords(ctx, db.GetCandidateRecordsParams{
		Limit:  int32(limit),
		Offset: int32(start),
	})
	if err != nil {
		return nil, err
	}
	recs := make([]*models.CandidateRecord, len(rows))
	for i, row := range rows {
		rec := &models.CandidateRecord{
			ID:          row.ID,
			SourceName:  row.SourceName,
			SourceID:    row.SourceID,
			Type:        row.Type,
			Metadata:    row.Metadata,
			DateCreated: row.DateCreated.Time,
			Status:      row.Status,
		}
		recs[i] = rec
	}
	return recs, err
}

func (r *Repo) CountCandidateRecords(ctx context.Context) (int, error) {
	num, err := r.queries.CountCandidateRecords(ctx)
	if err != nil {
		return 0, err
	}
	return int(num), nil
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

	return &models.CandidateRecord{
		ID:          row.ID,
		SourceName:  row.SourceName,
		SourceID:    row.SourceID,
		Type:        row.Type,
		Metadata:    row.Metadata,
		DateCreated: row.DateCreated.Time,
		Status:      row.Status,
	}, nil
}

func (r *Repo) GetCandidateRecord(ctx context.Context, id string) (*models.CandidateRecord, error) {
	row, err := r.queries.GetCandidateRecord(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &models.CandidateRecord{
		ID:          row.ID,
		SourceName:  row.SourceName,
		SourceID:    row.SourceID,
		Type:        row.Type,
		Metadata:    row.Metadata,
		DateCreated: row.DateCreated.Time,
		Status:      row.Status,
	}, nil
}

func (r *Repo) GetCandidateRecordAsPublication(ctx context.Context, id string) (*models.Publication, error) {
	cr, err := r.GetCandidateRecord(ctx, id)
	if err != nil {
		return nil, err
	}

	p := &models.Publication{}
	if err := json.Unmarshal(cr.Metadata, p); err != nil {
		return nil, err
	}
	for _, fn := range r.config.PublicationLoaders {
		if err := fn(p); err != nil {
			return nil, err
		}
	}
	return p, nil
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

func (r *Repo) ImportCandidateRecordAsPublication(ctx context.Context, rec *models.CandidateRecord, user *models.Person) (string, error) {
	var pubID string

	err := r.tx(ctx, func(r *Repo) error {
		pub := rec.AsPublication()
		pubID = ulid.Make().String()
		pub.ID = pubID
		if err := r.SavePublication(pub, user); err != nil {
			return err
		}

		if _, err := r.queries.SetCandidateRecordStatus(ctx, db.SetCandidateRecordStatusParams{ID: rec.ID, Status: "imported"}); err != nil {
			return err
		}

		return nil
	})

	return pubID, err
}
