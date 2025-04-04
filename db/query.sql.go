// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

const addCandidateRecord = `-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  id, source_name, source_id, source_metadata, type, metadata
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT(source_name, source_id)
DO
  UPDATE SET source_metadata = EXCLUDED.source_metadata, type = EXCLUDED.type, metadata = EXCLUDED.metadata
RETURNING id
`

type AddCandidateRecordParams struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Metadata       []byte
}

func (q *Queries) AddCandidateRecord(ctx context.Context, arg AddCandidateRecordParams) (string, error) {
	row := q.db.QueryRow(ctx, addCandidateRecord,
		arg.ID,
		arg.SourceName,
		arg.SourceID,
		arg.SourceMetadata,
		arg.Type,
		arg.Metadata,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const countPersonCandidateRecords = `-- name: CountPersonCandidateRecords :one
SELECT COUNT(*) FROM candidate_records WHERE status = 'new' AND (metadata->'author' @> $1::jsonb OR metadata->'supervisor' @> $1::jsonb)
`

func (q *Queries) CountPersonCandidateRecords(ctx context.Context, query []byte) (int64, error) {
	row := q.db.QueryRow(ctx, countPersonCandidateRecords, query)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCandidateRecord = `-- name: GetCandidateRecord :one
SELECT id, source_name, source_id, source_metadata, type, status, metadata, date_created, status_date, status_person_id, imported_id FROM candidate_records WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCandidateRecord(ctx context.Context, id string) (CandidateRecord, error) {
	row := q.db.QueryRow(ctx, getCandidateRecord, id)
	var i CandidateRecord
	err := row.Scan(
		&i.ID,
		&i.SourceName,
		&i.SourceID,
		&i.SourceMetadata,
		&i.Type,
		&i.Status,
		&i.Metadata,
		&i.DateCreated,
		&i.StatusDate,
		&i.StatusPersonID,
		&i.ImportedID,
	)
	return i, err
}

const getCandidateRecordBySource = `-- name: GetCandidateRecordBySource :one
SELECT id, source_name, source_id, source_metadata, type, status, metadata, date_created, status_date, status_person_id, imported_id FROM candidate_records WHERE source_name = $1 AND source_id = $2 LIMIT 1
`

type GetCandidateRecordBySourceParams struct {
	SourceName string
	SourceID   string
}

func (q *Queries) GetCandidateRecordBySource(ctx context.Context, arg GetCandidateRecordBySourceParams) (CandidateRecord, error) {
	row := q.db.QueryRow(ctx, getCandidateRecordBySource, arg.SourceName, arg.SourceID)
	var i CandidateRecord
	err := row.Scan(
		&i.ID,
		&i.SourceName,
		&i.SourceID,
		&i.SourceMetadata,
		&i.Type,
		&i.Status,
		&i.Metadata,
		&i.DateCreated,
		&i.StatusDate,
		&i.StatusPersonID,
		&i.ImportedID,
	)
	return i, err
}

const hasCandidateRecords = `-- name: HasCandidateRecords :one
SELECT EXISTS(SELECT 1 FROM candidate_records WHERE status = 'new')
`

func (q *Queries) HasCandidateRecords(ctx context.Context) (bool, error) {
	row := q.db.QueryRow(ctx, hasCandidateRecords)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const personHasCandidateRecords = `-- name: PersonHasCandidateRecords :one
SELECT EXISTS(SELECT 1 FROM candidate_records WHERE status = 'new' AND (metadata->'author' @> $1::jsonb OR metadata->'supervisor' @> $1::jsonb))
`

func (q *Queries) PersonHasCandidateRecords(ctx context.Context, query []byte) (bool, error) {
	row := q.db.QueryRow(ctx, personHasCandidateRecords, query)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const setCandidateRecordMetadata = `-- name: SetCandidateRecordMetadata :execresult
UPDATE candidate_records
SET metadata = $1
WHERE id = $2
`

type SetCandidateRecordMetadataParams struct {
	Metadata []byte
	ID       string
}

func (q *Queries) SetCandidateRecordMetadata(ctx context.Context, arg SetCandidateRecordMetadataParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, setCandidateRecordMetadata, arg.Metadata, arg.ID)
}

const setCandidateRecordStatus = `-- name: SetCandidateRecordStatus :one
UPDATE candidate_records 
SET status = $1,
    status_date = now(),
    status_person_id = $2,
    imported_id = $3
WHERE id = $4 RETURNING id
`

type SetCandidateRecordStatusParams struct {
	Status         string
	StatusPersonID *string
	ImportedID     *string
	ID             string
}

func (q *Queries) SetCandidateRecordStatus(ctx context.Context, arg SetCandidateRecordStatusParams) (string, error) {
	row := q.db.QueryRow(ctx, setCandidateRecordStatus,
		arg.Status,
		arg.StatusPersonID,
		arg.ImportedID,
		arg.ID,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}
