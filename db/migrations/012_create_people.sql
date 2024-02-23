CREATE TABLE people (
  id BIGSERIAL PRIMARY KEY,
  replaced_by_id BIGINT REFERENCES people (id) ON DELETE CASCADE,
  identifiers JSONB NOT NULL,
  name TEXT NOT NULL CHECK (name <> ''),
  preferred_name TEXT CHECK (preferred_name <> ''),
  given_name TEXT CHECK (given_name <> ''),
  preferred_given_name TEXT CHECK (preferred_given_name <> ''),
  family_name TEXT CHECK (family_name <> ''),
  preferred_family_name TEXT CHECK (preferred_family_name <> ''),
  honorific_prefix TEXT CHECK (honorific_prefix <> ''),
  email TEXT CHECK (email <> ''),
--   active BOOLEAN DEFAULT false NOT NULL,
--   username TEXT CHECK (username <> ''),
  attributes JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX people_updated_at_idx on people (updated_at);
CREATE INDEX people_identifiers_gin_idx on people using gin (identifiers);

---- create above / drop below ----

DROP TABLE PEOPLE CASCADE;
