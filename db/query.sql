-- name: AddCandidateRecord :one
INSERT INTO candidate_records (
  source_name, source_id, metadata
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: GetCandidateRecordBySource :one
SELECT * FROM candidate_records
WHERE source_name = $1 AND source_id = $2;