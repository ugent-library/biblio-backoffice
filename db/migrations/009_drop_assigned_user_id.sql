ALTER TABLE candidate_records DROP COLUMN assigned_user_id;

---- create above / drop below ----

ALTER TABLE candidate_records ADD COLUMN assigned_user_id TEXT;