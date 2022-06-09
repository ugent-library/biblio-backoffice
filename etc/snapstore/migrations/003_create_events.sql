create table events (
    id text primary key,
    seq bigserial not null,
    stream_id text not null,
    stream_type text not null,
    type text not null,
    data jsonb,
    meta jsonb,
    date_created timestamptz not null default now()
);

create table snapshots (
    stream_id text primary key,
    stream_type text not null,
    event_id text not null,
    data jsonb not null,
    date_created timestamptz not null,
    date_updated timestamptz not null
);

---- create above / drop below ----

drop table events cascade;
drop table snapshots cascade;
