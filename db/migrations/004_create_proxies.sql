create table proxies (
    person_id text not null,
    proxy_person_id text not null,
    date_created timestamptz not null default now(),
    unique (person_id, proxy_person_id)
);

create index proxies_person_id_idx on proxies(person_id);
create index proxies_proxy_person_id_idx on proxies(proxy_person_id);

---- create above / drop below ----

drop table proxies cascade;
