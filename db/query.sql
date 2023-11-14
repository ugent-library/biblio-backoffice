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

-- name: GetCandidateRecordBySource :one
SELECT * FROM candidate_records
WHERE source_name = $1 AND source_id = $2;