create table projects (
  id bigserial primary key,
  name jsonb null,
  description jsonb null,
  founding_date text null,
  dissolution_date text null,
  attributes jsonb,
  identifiers jsonb,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp
);

create index projects_updated_at_idx on projects (updated_at);

---- create above / drop below ----

drop table projects CASCADE;