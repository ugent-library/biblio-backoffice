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
