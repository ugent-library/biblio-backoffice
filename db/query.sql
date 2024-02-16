-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  id, source_name, source_id, source_metadata, type, metadata
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT(source_name, source_id)
DO
  UPDATE SET source_metadata = EXCLUDED.source_metadata, type = EXCLUDED.type, metadata = EXCLUDED.metadata
RETURNING id;

-- name: GetCandidateRecords :many
SELECT * FROM candidate_records WHERE status = 'new' ORDER BY date_created ASC LIMIT $1 OFFSET $2;

-- name: CountCandidateRecords :one
SELECT count(*) count FROM candidate_records WHERE status = 'new';

-- name: GetCandidateRecord :one
SELECT * FROM candidate_records WHERE status = 'new' AND id = $1 LIMIT 1;

-- name: SetStatusCandidateRecord :one
UPDATE candidate_records SET status = $1 WHERE id = $2 RETURNING id;

-- name: GetCandidateRecordBySource :one
SELECT * FROM candidate_records WHERE source_name = $1 AND source_id = $2 LIMIT 1;