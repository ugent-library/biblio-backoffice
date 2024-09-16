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
SELECT *, count(*) OVER () AS total
FROM candidate_records
WHERE (status = 'new' OR  EXTRACT(DAY FROM (current_timestamp - status_date)) <= 90)
ORDER BY date_created ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: GetCandidateRecordsByPersonID :many
SELECT *, count(*) OVER () AS total
FROM candidate_records
WHERE (status = 'new' OR (sqlc.arg('new_only')::bool = 0::bool AND EXTRACT(DAY FROM (current_timestamp - status_date)) <= 90))
  AND (metadata->'author' @> sqlc.arg(query)::jsonb OR metadata->'supervisor' @> sqlc.arg(query)::jsonb)
ORDER BY date_created ASC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: HasCandidateRecords :one
SELECT EXISTS(SELECT 1 FROM candidate_records WHERE status = 'new');

-- name: PersonHasCandidateRecords :one
SELECT EXISTS(SELECT 1 FROM candidate_records WHERE status = 'new' AND (metadata->'author' @> sqlc.arg(query)::jsonb OR metadata->'supervisor' @> sqlc.arg(query)::jsonb));


-- name: GetCandidateRecord :one
SELECT * FROM candidate_records WHERE id = $1 LIMIT 1;

-- name: SetCandidateRecordStatus :one
UPDATE candidate_records 
SET status = sqlc.arg('status'),
    status_date = now(),
    status_person_id = sqlc.arg('status_person_id'),
    imported_id = sqlc.arg('imported_id')
WHERE id = sqlc.arg('id') RETURNING id;

-- name: GetCandidateRecordBySource :one
SELECT * FROM candidate_records WHERE source_name = $1 AND source_id = $2 LIMIT 1;