create extension if not exists btree_gist;

create table mutations (
    mutation_id text primary key,
    seq bigserial not null,
    entity_id text not null,
    entity_type text not null,
    exclude using gist (entity_id with =, entity_type with <>),
    mutation_name text not null,
    mutation_data jsonb,
    mutation_meta jsonb,
    date_created timestamptz not null default now()
);

create table projections (
    entity_id text,
    entity_type text,
    primary key (entity_id, entity_type),
    entity_data jsonb not null,
    mutation_id text not null references mutations(mutation_id),
    date_created timestamptz not null,
    date_updated timestamptz not null
);

---- create above / drop below ----

drop table projections cascade;
drop table mutations cascade;
