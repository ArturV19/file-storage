DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'file_storage') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE file_storage');
END IF;
END
$$;

\c file_storage;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
