create table datasets (
    snapshot_id text primary key,
    id text not null,
    -- snapshot_id uuid primary key default gen_random_uuid(),
    -- id uuid not null,
    data jsonb not null,
    date_from timestamptz not null default now(),
    date_until timestamptz not null default 'infinity'::timestamptz
);

create index datasets_snapshot_id_idx on datasets(snapshot_id);
create index datasets_id_idx on datasets(id);
create index datasets_date_from_idx on datasets(date_from);
create index datasets_date_until_idx on datasets(date_until);

---- create above / drop below ----

drop table datasets cascade;
