-- Creates rw role
-- the default role will be a superuser (rolsuper = true) and has the BYPASSRLS privilege (rolbypassrls = true).
-- which have the ability to bypass row-level security (RLS) policies.
CREATE ROLE rw WITH LOGIN PASSWORD 'password' NOSUPERUSER NOCREATEDB NOCREATEROLE;

-- Grant all privileges on future tables to rw
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO rw;
