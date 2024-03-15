ALTER TABLE people ADD COLUMN publication_count INT NOT NULL DEFAULT 0;
ALTER TABLE projects ADD COLUMN publication_count INT NOT NULL DEFAULT 0;

---- create above / drop below ----

ALTER TABLE people DROP COLUMN publication_count INT;
ALTER TABLE projects DROP COLUMN publication_count INT;
