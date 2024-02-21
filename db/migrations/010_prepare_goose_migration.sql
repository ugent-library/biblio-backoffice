-- This migration only exists to prepare a switch to the goose migration tool.
-- The last version_id should match the desired number of goose migrations.

CREATE TABLE goose_migration (
    id SERIAL PRIMARY KEY,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);

INSERT INTO goose_migration (version_id, is_applied) VALUES (0, TRUE);
INSERT INTO goose_migration (version_id, is_applied) VALUES (1, TRUE);

---- create above / drop below ----

DROP TABLE goose_migration CASCADE;