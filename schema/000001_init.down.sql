DO $$
    DECLARE
        v_table_name text;
    BEGIN
        FOR v_table_name IN (SELECT table_name FROM information_schema.tables WHERE table_name LIKE 'messages_%')
            LOOP
                EXECUTE 'DROP TABLE IF EXISTS ' || v_table_name || ' CASCADE';
            END LOOP;
    END $$;

DROP TABLE rassilka;

DROP TABLE clients;

DROP SEQUENCE IF EXISTS messages_id_seq;
