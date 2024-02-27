create table projects (
  id bigserial primary key,
  replaced_by_id bigint references projects (id) on delete cascade,
  names jsonb null,
  descriptions jsonb null,
  founding_date text null,
  dissolution_date text null,
  deleted boolean not null default false,
  attributes jsonb,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp
);

create index projects_updated_at_idx on projects (updated_at);
create index projects_replaced_by_id_fkey on projects (replaced_by_id);

create table project_identifiers (
  project_id bigint not null references projects on delete cascade,
  kind text not null check (kind <> ''),
  value text not null check (value <> ''),
  unique (project_id, kind, value)
);

create index project_identifiers_project_id_fkey on project_identifiers (project_id);
create index project_identifiers_kind_value_idx on project_identifiers (kind, value);

---- create above / drop below ----

drop table project_identifiers cascade;
drop table projects cascade;