CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Creates rw role
-- the default role will be a superuser (rolsuper = true) and has the BYPASSRLS privilege (rolbypassrls = true).
-- which have the ability to bypass row-level security (RLS) policies.
CREATE ROLE rw WITH LOGIN PASSWORD 'password' NOSUPERUSER NOCREATEDB NOCREATEROLE;

-- Grant all privileges on future tables to rw
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO rw;

-- Create tables
CREATE TABLE users (
  slug citext PRIMARY KEY,
  description TEXT,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE todos (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  description TEXT,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  user_slug CITEXT NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user_slug FOREIGN KEY (user_slug) REFERENCES users (slug) 
    ON UPDATE CASCADE ON DELETE CASCADE
);

-- Row Level Security (RLS) Policy
CREATE POLICY user_isolation_policy ON users
  USING (slug = current_setting('app.current_user'));

CREATE POLICY todo_isolation_policy ON todos
  USING (user_slug = current_setting('app.current_user'));

-- Enable RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE todos ENABLE ROW LEVEL SECURITY;

-- Force RLS on table owner
-- by default table owner bypasses row-level security
ALTER TABLE users FORCE ROW LEVEL SECURITY;
ALTER TABLE todos FORCE ROW LEVEL SECURITY;

-- Seed users
INSERT INTO users (slug, description)
VALUES
  ('user-1', 'default user 1 description'),
  ('user-2', 'default user 2 description'),
  ('user-3', 'default user 3 description')
ON CONFLICT DO NOTHING;
