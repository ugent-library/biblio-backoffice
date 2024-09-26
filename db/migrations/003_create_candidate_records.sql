create table candidate_records (
    id text primary key,
    source_name text not null,
    source_id text not null,
    source_metadata bytea not null,
    type text not null,
    status text not null default 'new',
    metadata jsonb not null,
    date_created timestamptz not null default now(),
    unique(source_name, source_id)
);

create index candidate_records_source_idx on candidate_records (source_name, source_id);

---- create above / drop below ----

drop table candidate_records cascade;
