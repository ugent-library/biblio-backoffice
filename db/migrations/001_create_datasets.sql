create table datasets (
    snapshot_id text primary key,
    id text not null,
    data jsonb not null,
    date_from timestamptz not null default now(),
    date_until timestamptz
);

create index datasets_snapshot_id_idx on datasets(snapshot_id);
create index datasets_id_idx on datasets(id);
create index datasets_date_from_idx on datasets(date_from);
create index datasets_date_until_idx on datasets(date_until);
create index if not exists datasets_date_created_idx on datasets((data->>'date_created'));
create index if not exists datasets_date_updated_idx on datasets((data->>'date_updated'));

---- create above / drop below ----

drop table datasets cascade;
