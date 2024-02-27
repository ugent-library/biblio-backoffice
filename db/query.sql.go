// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package db

import (
	"context"
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

const countCandidateRecords = `-- name: CountCandidateRecords :one
SELECT count(*) count FROM candidate_records WHERE status = 'new'
`

func (q *Queries) CountCandidateRecords(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCandidateRecords)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCandidateRecord = `-- name: GetCandidateRecord :one
SELECT id, source_name, source_id, source_metadata, type, metadata, date_created, status FROM candidate_records WHERE status = 'new' AND id = $1 LIMIT 1
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
		&i.Metadata,
		&i.DateCreated,
		&i.Status,
	)
	return i, err
}

const getCandidateRecordBySource = `-- name: GetCandidateRecordBySource :one
SELECT id, source_name, source_id, source_metadata, type, metadata, date_created, status FROM candidate_records WHERE source_name = $1 AND source_id = $2 LIMIT 1
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
		&i.Metadata,
		&i.DateCreated,
		&i.Status,
	)
	return i, err
}

const getCandidateRecords = `-- name: GetCandidateRecords :many
SELECT id, source_name, source_id, source_metadata, type, metadata, date_created, status FROM candidate_records WHERE status = 'new' ORDER BY date_created ASC LIMIT $1 OFFSET $2
`

type GetCandidateRecordsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetCandidateRecords(ctx context.Context, arg GetCandidateRecordsParams) ([]CandidateRecord, error) {
	rows, err := q.db.Query(ctx, getCandidateRecords, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CandidateRecord
	for rows.Next() {
		var i CandidateRecord
		if err := rows.Scan(
			&i.ID,
			&i.SourceName,
			&i.SourceID,
			&i.SourceMetadata,
			&i.Type,
			&i.Metadata,
			&i.DateCreated,
			&i.Status,
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

const setCandidateRecordStatus = `-- name: SetCandidateRecordStatus :one
UPDATE candidate_records SET status = $1 WHERE id = $2 RETURNING id
`

type SetCandidateRecordStatusParams struct {
	Status string
	ID     string
}

func (q *Queries) SetCandidateRecordStatus(ctx context.Context, arg SetCandidateRecordStatusParams) (string, error) {
	row := q.db.QueryRow(ctx, setCandidateRecordStatus, arg.Status, arg.ID)
	var id string
	err := row.Scan(&id)
	return id, err
}
