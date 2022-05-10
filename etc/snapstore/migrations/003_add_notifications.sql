DO $$
BEGIN
    CREATE OR REPLACE FUNCTION publication_snapshots_notify() RETURNS trigger AS $FN$
    DECLARE
        payload jsonb;
    BEGIN
        payload = jsonb_build_object(
            'id', to_jsonb(NEW.id),
            'data', to_jsonb(NEW.data)
        );

        PERFORM pg_notify('publication_snapshots', payload::TEXT);

        RETURN NEW;
    END;
    $FN$ LANGUAGE plpgsql;

	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'publication_snapshots_insert') THEN
		CREATE TRIGGER publication_snapshots_insert AFTER INSERT ON publication_snapshots FOR EACH ROW EXECUTE PROCEDURE publication_snapshots_notify();
	END IF;
END $$;

---- create above / drop below ----
