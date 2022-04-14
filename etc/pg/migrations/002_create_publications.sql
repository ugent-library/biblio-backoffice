create table publications (
    pk bigserial primary key,
    id text not null,
    data jsonb not null,
    data_from timestamp with time zone not null,
    data_to timestamp with time zone
);

create index publications_id_idx on publications(id);
create index publications_data_from_idx on publications(data_from);
create index publications_data_to_idx on publications(data_to);

---- create above / drop below ----

drop table publications cascade;
