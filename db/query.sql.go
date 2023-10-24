// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package db

import (
	"context"
)

const addCandidateRecord = `-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  source_name, source_id, metadata
) VALUES (
  $1, $2, $3
)
RETURNING id
`

type AddCandidateRecordParams struct {
	SourceName string
	SourceID   string
	Metadata   []byte
}

func (q *Queries) AddCandidateRecord(ctx context.Context, arg AddCandidateRecordParams) (string, error) {
	row := q.db.QueryRow(ctx, addCandidateRecord, arg.SourceName, arg.SourceID, arg.Metadata)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getCandidateRecordBySource = `-- name: GetCandidateRecordBySource :one
SELECT id, metadata, source_name, source_id, date_created FROM candidate_records
WHERE source_name = $1 AND source_id = $2
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
		&i.Metadata,
		&i.SourceName,
		&i.SourceID,
		&i.DateCreated,
	)
	return i, err
}
