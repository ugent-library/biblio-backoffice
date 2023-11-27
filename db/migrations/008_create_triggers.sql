create function notify_record() returns trigger as $$
    declare
        record_type text := TG_ARGV[0];
        evt text;
    begin
        if (TG_OP = 'DELETE') then
            evt = json_build_object(
                'name', 'purge',
                'record_type', record_type,
                'record_id', OLD.id
            )::text;
        elsif (TG_OP = 'INSERT') then 
            evt = json_build_object(
                'name', 'create',
                'record_type', record_type,
                'record_id', NEW.id,
                'snapshot_id', NEW.snapshot_id
            )::text;
        elsif (TG_OP = 'UPDATE') then 
            evt = json_build_object(
                'name', 'update',
                'record_type', record_type,
                'record_id', NEW.id,
                'snapshot_id', NEW.snapshot_id
            )::text;
        end if;

        perform pg_notify('events', evt);

        return null;
    end;
$$ language plpgsql;

create trigger notify_dataset_changed
after insert or update on datasets
    for each row
    when (NEW.date_until is null)
    execute procedure notify_record('Dataset');

create trigger notify_dataset_deleted
after delete on datasets
    for each row
    when (OLD.date_until is null)
    execute procedure notify_record('Dataset');

create trigger notify_publication_changed
after insert on publications
    for each row
    when (NEW.date_until is null)
    execute procedure notify_record('Publication');

create trigger notify_publication_deleted
after delete on publications
    for each row
    when (OLD.date_until is null)
    execute procedure notify_record('Publication');

---- create above / drop below ----

drop trigger notify_dataset_changed on datasets;
drop trigger notify_dataset_deleted on datasets;
drop trigger notify_publication_changed on publications;
drop trigger notify_publication_deleted on publications;
drop function notify_record();