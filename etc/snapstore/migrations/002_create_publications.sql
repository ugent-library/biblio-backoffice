create table publication_snapshots (
    snapshot_id bigserial primary key,
    id text not null,
    data jsonb not null,
    date_from timestamptz not null default now(),
    date_until timestamptz not null default 'infinity'::timestamptz
);

create table publication_versions (
    version_id bigserial primary key,
    affinity_id text not null,
    snapshot_id bigint,
    id text not null,
    data jsonb not null,
    date_created timestamptz not null default now(),
    constraint snapshot_fkey foreign key(snapshot_id) references publication_snapshots(snapshot_id)
);

create index publication_snapshots_id_idx on publication_snapshots(id);
create index publication_snapshots_date_from_idx on publication_snapshots(date_from);
create index publication_snapshots_date_until_idx on publication_snapshots(date_until);

create index publication_versions_id_idx on publication_versions(affinity_id, id);
create index publication_versions_date_created_idx on publication_versions(date_created);

---- create above / drop below ----

drop table publication_snapshots cascade;
drop table publication_versions cascade;
