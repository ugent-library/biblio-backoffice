ALTER TABLE organizations ADD COLUMN position INT;
CREATE INDEX organizations_position_idx ON organizations (position);

---- create above / drop below ----

DROP INDEX organizations_position_idx;
ALTER TABLE organizations DROP COLUMN position;
