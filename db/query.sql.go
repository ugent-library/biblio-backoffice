// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
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

const getCandidateRecords = `-- name: GetCandidateRecords :many
SELECT id, source_name, source_id, source_metadata, type, status, metadata, date_created, status_date, status_person_id, imported_id, count(*) OVER () AS total
FROM candidate_records
WHERE (status = 'new' OR (status_date IS NOT NULL AND EXTRACT(DAY FROM (current_timestamp - status_date)) <= 90))
ORDER BY date_created ASC
LIMIT $2
OFFSET $1
`

type GetCandidateRecordsParams struct {
	Offset int32
	Limit  int32
}

type GetCandidateRecordsRow struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Status         string
	Metadata       []byte
	DateCreated    pgtype.Timestamptz
	StatusDate     pgtype.Timestamptz
	StatusPersonID *string
	ImportedID     *string
	Total          int64
}

func (q *Queries) GetCandidateRecords(ctx context.Context, arg GetCandidateRecordsParams) ([]GetCandidateRecordsRow, error) {
	rows, err := q.db.Query(ctx, getCandidateRecords, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCandidateRecordsRow
	for rows.Next() {
		var i GetCandidateRecordsRow
		if err := rows.Scan(
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
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCandidateRecordsByPersonID = `-- name: GetCandidateRecordsByPersonID :many
SELECT id, source_name, source_id, source_metadata, type, status, metadata, date_created, status_date, status_person_id, imported_id, count(*) OVER () AS total
FROM candidate_records
WHERE (status = 'new' OR (status_date IS NOT NULL AND $1::bool = 0::bool AND EXTRACT(DAY FROM (current_timestamp - status_date)) <= 90))
  AND (metadata->'author' @> $2::jsonb OR metadata->'supervisor' @> $2::jsonb)
ORDER BY date_created ASC
LIMIT $4
OFFSET $3
`

type GetCandidateRecordsByPersonIDParams struct {
	NewOnly bool
	Query   []byte
	Offset  int32
	Limit   int32
}

type GetCandidateRecordsByPersonIDRow struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Status         string
	Metadata       []byte
	DateCreated    pgtype.Timestamptz
	StatusDate     pgtype.Timestamptz
	StatusPersonID *string
	ImportedID     *string
	Total          int64
}

func (q *Queries) GetCandidateRecordsByPersonID(ctx context.Context, arg GetCandidateRecordsByPersonIDParams) ([]GetCandidateRecordsByPersonIDRow, error) {
	rows, err := q.db.Query(ctx, getCandidateRecordsByPersonID,
		arg.NewOnly,
		arg.Query,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCandidateRecordsByPersonIDRow
	for rows.Next() {
		var i GetCandidateRecordsByPersonIDRow
		if err := rows.Scan(
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
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
