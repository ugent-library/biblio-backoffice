// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addCandidateRecord = `-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  id, source_name, source_id, source_metadata, type, metadata, assigned_user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT(source_name, source_id)
DO
  UPDATE SET source_metadata = EXCLUDED.source_metadata, type = EXCLUDED.type, metadata = EXCLUDED.metadata, assigned_user_id = EXCLUDED.assigned_user_id
RETURNING id
`

type AddCandidateRecordParams struct {
	ID             string
	SourceName     string
	SourceID       string
	SourceMetadata []byte
	Type           string
	Metadata       []byte
	AssignedUserID pgtype.Text
}

func (q *Queries) AddCandidateRecord(ctx context.Context, arg AddCandidateRecordParams) (string, error) {
	row := q.db.QueryRow(ctx, addCandidateRecord,
		arg.ID,
		arg.SourceName,
		arg.SourceID,
		arg.SourceMetadata,
		arg.Type,
		arg.Metadata,
		arg.AssignedUserID,
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
SELECT id, source_name, source_id, source_metadata, type, metadata, assigned_user_id, date_created, status FROM candidate_records WHERE status = 'new' AND id = $1 LIMIT 1
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
		&i.AssignedUserID,
		&i.DateCreated,
		&i.Status,
	)
	return i, err
}

const getCandidateRecords = `-- name: GetCandidateRecords :many
SELECT id, source_name, source_id, source_metadata, type, metadata, assigned_user_id, date_created, status FROM candidate_records WHERE status = 'new' ORDER BY date_created ASC LIMIT $1 OFFSET $2
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
			&i.AssignedUserID,
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

const getCandidateRecordsByUser = `-- name: GetCandidateRecordsByUser :many
SELECT id, source_name, source_id, source_metadata, type, metadata, assigned_user_id, date_created, status FROM candidate_records
WHERE assigned_user_id = $1 AND status = 'new'
`

func (q *Queries) GetCandidateRecordsByUser(ctx context.Context, assignedUserID pgtype.Text) ([]CandidateRecord, error) {
	rows, err := q.db.Query(ctx, getCandidateRecordsByUser, assignedUserID)
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
			&i.AssignedUserID,
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

const setStatusCandidateRecord = `-- name: SetStatusCandidateRecord :one
UPDATE candidate_records SET status = $1 WHERE id = $2 RETURNING id
`

type SetStatusCandidateRecordParams struct {
	Status string
	ID     string
}

func (q *Queries) SetStatusCandidateRecord(ctx context.Context, arg SetStatusCandidateRecordParams) (string, error) {
	row := q.db.QueryRow(ctx, setStatusCandidateRecord, arg.Status, arg.ID)
	var id string
	err := row.Scan(&id)
	return id, err
}
