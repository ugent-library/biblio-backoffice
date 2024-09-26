create table proxies (
    proxy_person_id text not null check (person_id <> ''),
    person_id text not null check (person_id <> ''),
    date_created timestamptz not null default now(),
    unique (proxy_person_id, person_id),
    check (proxy_person_id <> person_id)
);

create index proxies_proxy_person_id_idx on proxies(proxy_person_id);
create index proxies_person_id_idx on proxies(person_id);

---- create above / drop below ----

drop table proxies cascade;
