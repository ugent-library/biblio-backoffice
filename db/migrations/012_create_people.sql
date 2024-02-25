CREATE TABLE people (
  id BIGSERIAL PRIMARY KEY,
  replaced_by_id BIGINT REFERENCES people (id) ON DELETE CASCADE,
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

CREATE INDEX people_replaced_by_id_fkey on people (replaced_by_id);
CREATE INDEX people_updated_at_idx on people (updated_at);

CREATE TABLE person_identifiers (
  person_id BIGINT NOT NULL REFERENCES people ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind <> ''),
  value TEXT NOT NULL CHECK (value <> ''),
  UNIQUE (person_id, kind, value)
);

CREATE INDEX person_identifiers_person_id_fkey ON person_identifiers (person_id);
CREATE INDEX person_identifiers_kind_value_idx ON person_identifiers (kind, value);

---- create above / drop below ----

DROP TABLE person_identifiers CASCADE;
DROP TABLE people CASCADE;
