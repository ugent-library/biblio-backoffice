create index candidate_records_date_created_idx on candidate_records (date_created);
create index candidate_records_date_author_idx on candidate_records using gin((metadata->'author') jsonb_path_ops);
create index candidate_records_date_supervisor_idx on candidate_records using gin((metadata->'supervisor') jsonb_path_ops);

---- create above / drop below ----

drop index candidate_records_date_created_idx;
drop index candidate_records_date_author_idx;
drop index candidate_records_date_supervisor_idx;
