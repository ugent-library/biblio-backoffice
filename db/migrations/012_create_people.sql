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
  active BOOLEAN NOT NULL DEFAULT false,
  role TEXT CHECK (role <> ''),
  username TEXT CHECK (username <> ''),
  attributes JSONB,
  tokens JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX people_replaced_by_id_fkey on people (replaced_by_id);
CREATE INDEX people_active_idx ON people (active) WHERE (replaced_by_id IS NULL);
CREATE INDEX people_username_idx ON people (username) WHERE (replaced_by_id IS NULL);
CREATE INDEX people_updated_at_idx on people (updated_at) WHERE (replaced_by_id IS NULL);

CREATE TABLE person_identifiers (
  person_id BIGINT NOT NULL REFERENCES people ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind <> ''),
  value TEXT NOT NULL CHECK (value <> '')
);

CREATE INDEX person_identifiers_person_id_fkey ON person_identifiers (person_id);
CREATE INDEX person_identifiers_kind_value_idx ON person_identifiers (kind, value);

CREATE TABLE organizations (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT REFERENCES organizations ON DELETE SET NULL CHECK (parent_id <> id),
  names JSONB,
  ceased BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_identifiers (
  organization_id BIGINT NOT NULL REFERENCES organizations ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind <> ''),
  value TEXT NOT NULL CHECK (value <> '')
);

CREATE INDEX organization_identifiers_organization_id_fkey ON organization_identifiers (organization_id);
CREATE INDEX organization_identifiers_kind_value_idx ON organization_identifiers (kind, value);

CREATE TABLE affiliations (
  person_id BIGINT NOT NULL REFERENCES people ON DELETE CASCADE,
  organization_id BIGINT NOT NULL REFERENCES organizations ON DELETE CASCADE,
  UNIQUE (person_id, organization_id)
);

---- create above / drop below ----

DROP TABLE affiliations CASCADE;

DROP TABLE person_identifiers CASCADE;
DROP TABLE people CASCADE;

DROP TABLE organization_identifiers CASCADE;
DROP TABLE organizations CASCADE;