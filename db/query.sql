-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  id, source_name, source_id, source_metadata, type, metadata, assigned_user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT(source_name, source_id)
DO
  UPDATE SET source_metadata = EXCLUDED.source_metadata, type = EXCLUDED.type, metadata = EXCLUDED.metadata, assigned_user_id = EXCLUDED.assigned_user_id
RETURNING id;

-- name: GetCandidateRecordsByUser :many
SELECT * FROM candidate_records
WHERE assigned_user_id = $1 AND status = 'new';

-- name: GetCandidateRecords :many
SELECT * FROM candidate_records WHERE status = 'new' ORDER BY date_created ASC LIMIT $1 OFFSET $2;

-- name: CountCandidateRecords :one
SELECT count(*) count FROM candidate_records WHERE status = 'new';

-- name: GetCandidateRecord :one
SELECT * FROM candidate_records WHERE status = 'new' AND id = $1 LIMIT 1;

-- name: SetStatusCandidateRecord :one
UPDATE candidate_records SET status = $1 WHERE id = $2 RETURNING id;
