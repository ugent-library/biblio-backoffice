create table dataset_snapshots (
    snapshot_id bigserial primary key,
    id text not null,
    data jsonb not null,
    date_from timestamptz not null default now(),
    date_until timestamptz not null default 'infinity'::timestamptz
);

create table dataset_versions (
    version_id bigserial primary key,
    affinity_id text not null,
    snapshot_id bigint,
    id text not null,
    data jsonb not null,
    date_created timestamptz not null default now(),
    constraint snapshot_fkey foreign key(snapshot_id) references dataset_snapshots(snapshot_id)
);

create index dataset_snapshots_id_idx on dataset_snapshots(id);
create index dataset_snapshots_date_from_idx on dataset_snapshots(date_from);
create index dataset_snapshots_date_until_idx on dataset_snapshots(date_until);

create index dataset_versions_id_idx on dataset_versions(affinity_id, id);
create index dataset_versions_date_created_idx on dataset_versions(date_created);

---- create above / drop below ----

drop table dataset_snapshots cascade;
drop table dataset_versions cascade;
