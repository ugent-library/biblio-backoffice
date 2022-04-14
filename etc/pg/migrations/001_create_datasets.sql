create table datasets (
    pk bigserial primary key,
    id text not null,
    data jsonb not null,
    data_from timestamp with time zone not null,
    data_to timestamp with time zone
);

create index datasets_id_idx on datasets(id);
create index datasets_data_from_idx on datasets(data_from);
create index datasets_data_to_idx on datasets(data_to);

---- create above / drop below ----

drop table datasets cascade;
