-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  id, source_name, source_id, source_metadata, type, metadata
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: GetCandidateRecordBySource :one
SELECT * FROM candidate_records
WHERE source_name = $1 AND source_id = $2;