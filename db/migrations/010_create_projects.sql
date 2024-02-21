create table projects (
  id bigserial primary key,
  names jsonb null,
  descriptions jsonb null,
  founding_date text null,
  dissolution_date text null,
  attributes jsonb,
  identifiers jsonb,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp
);

create index projects_updated_at_idx on projects (updated_at);
create index projects_identifiers_gin_idx on projects using gin (identifiers);

---- create above / drop below ----

drop table projects CASCADE;