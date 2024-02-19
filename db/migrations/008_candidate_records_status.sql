alter table candidate_records add column status text not null default 'new';

---- create above / drop below ----
alter table candidate_records drop column status;
