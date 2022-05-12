create table publications (
    snapshot_id text primary key,
    id text not null,
    -- snapshot_id uuid primary key default gen_random_uuid(),
    -- id uuid not null,
    data jsonb not null,
    date_from timestamptz not null default now(),
    date_until timestamptz not null default 'infinity'::timestamptz
);

create index publications_snapshot_id_idx on publications(snapshot_id);
create index publications_id_idx on publications(id);
create index publications_date_from_idx on publications(date_from);
create index publications_date_until_idx on publications(date_until);

---- create above / drop below ----

drop table publications cascade;
