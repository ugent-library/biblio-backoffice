create extension if not exists btree_gist;

create table events (
    id text primary key,
    seq bigserial not null,
    stream_id text not null,
    stream_type text not null,
    exclude using gist (stream_id with =, stream_type with <>),
    name text not null,
    data jsonb,
    meta jsonb,
    date_created timestamptz not null default now()
);

create table projections (
    stream_id text,
    stream_type text,
    primary key (stream_id, stream_type),
    event_id text not null references events(id),
    data jsonb not null,
    date_created timestamptz not null,
    date_updated timestamptz not null
);

---- create above / drop below ----

drop table projections cascade;
drop table events cascade;
drop extension btree_gist;
