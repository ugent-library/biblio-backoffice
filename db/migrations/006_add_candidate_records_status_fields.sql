alter table candidate_records
    add column status_date timestamptz null,
    add column status_person_id text null,
    add column imported_id text null;

---- create above / drop below ----

alter table candidate_records
    drop column status_date,
    drop column status_person_id,
    drop column imported_id;