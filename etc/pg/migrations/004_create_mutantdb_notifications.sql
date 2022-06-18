create function notify_projection_changed() returns trigger as $$
    declare
        payload_id text;
        payload text;
        size int;
        n int = 0;
    begin
        payload_id = NEW.mutation_id;
        payload = row_to_json(NEW)::text;
        size = length(payload);

        -- send payload in chunks of 4000 bytes
        while n <= size loop
            -- append counter to ensure non-identical payload strings
            perform pg_notify(
                'projections',
                concat(payload_id, n::text, ':', substr(payload, n, 4000))
            );
            n = n + 4000;
	    end loop;

        perform pg_notify('projections', concat(payload_id, 'EOF'));

        return null;
    end;
$$ language plpgsql;

create trigger notify_projection_insert_update
after insert or update on projections
    for each row execute procedure notify_projection_changed();

---- create above / drop below ----

drop trigger notify_projection_insert_update on projections;
drop function notify_projection_changed();