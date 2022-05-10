do $$
begin
    create function publication_snapshots_notify()
    returns trigger as $fn$
    declare
        payload jsonb;
    begin
        payload = jsonb_build_object(
            'id', to_jsonb(new.id),
            'data', to_jsonb(new.data),
            'date_from', to_jsonb(new.date_from),
            'date_until', to_jsonb(new.date_until)
        );

        perform pg_notify('publication:snapshot', payload::text);

        return new;
    end;
    $fn$
    language plpgsql;

    create trigger publication_snapshots_insert
    after insert on publication_snapshots
    for each row execute procedure publication_snapshots_notify();
end $$;

---- create above / drop below ----

drop trigger publication_snapshots_insert on publication_snapshots;
drop function publication_snapshots_notify();
