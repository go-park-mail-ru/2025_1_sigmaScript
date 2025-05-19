DO $$ 
BEGIN
   IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'schema_editor_role') THEN
      DROP ROLE schema_editor_role;
   END IF;
END $$;

DO $$ BEGIN
   IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'service_user') THEN
      DROP USER service_user;
   END IF;
END $$;


CREATE ROLE schema_editor_role;

GRANT USAGE ON SCHEMA public TO schema_editor_role;
GRANT CREATE ON SCHEMA public TO schema_editor_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO schema_editor_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, UPDATE, DELETE ON TABLES TO schema_editor_role;

CREATE USER service_user WITH PASSWORD 'service_user_password';

GRANT schema_editor_role TO service_user;
