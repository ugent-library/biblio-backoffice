package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/db"
	"github.com/ugent-library/biblio-backoffice/models"
)

func (r *Repo) AddCandidateRecord(ctx context.Context, rec *models.CandidateRecord) error {
	params := db.AddCandidateRecordParams{
		ID:             ulid.Make().String(),
		SourceName:     rec.SourceName,
		SourceID:       rec.SourceID,
		SourceMetadata: rec.SourceMetadata,
		Type:           rec.Type,
		Metadata:       rec.Metadata,
		AssignedUserID: pgtype.Text{String: rec.AssignedUserID, Valid: rec.AssignedUserID != ""},
	}
	_, err := r.queries.AddCandidateRecord(ctx, params)
	return err
}

func (r *Repo) GetCandidateRecordsByUser(ctx context.Context, userID string) ([]*models.CandidateRecord, error) {
	rows, err := r.queries.GetCandidateRecordsByUser(ctx, pgtype.Text{String: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	recs := make([]*models.CandidateRecord, len(rows))
	for i, row := range rows {
		rec := &models.CandidateRecord{
			SourceName:     row.SourceName,
			SourceID:       row.SourceID,
			AssignedUserID: row.AssignedUserID.String,
			DateCreated:    row.DateCreated.Time,
		}
		recs[i] = rec
	}
	return recs, err
}
