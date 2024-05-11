-- Ensure that the citext and pgcrypto extensions are available in the database.
CREATE EXTENSION IF NOT EXISTS citext;

-- Create a role `rw` with login permission and no superuser or RLS-bypass privileges.
-- This role will be used by the app, as configured in the `.env` file.
CREATE ROLE rw WITH LOGIN PASSWORD 'password' NOSUPERUSER NOCREATEDB NOCREATEROLE;

-- Grant the `rw` role all privileges on tables created in the `public` schema.
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO rw;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO rw;

-- Create the `users` table:
CREATE TABLE users (
  slug citext PRIMARY KEY,
  description TEXT,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Create the `todos` table:
CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  user_slug CITEXT NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user_slug FOREIGN KEY (user_slug) REFERENCES users (slug) 
    ON UPDATE CASCADE ON DELETE CASCADE
);

-- Define Row-Level Security (RLS) policies:

-- Access to rows in the `users` table is restricted by a policy that compares the `slug` column with the `current session's tenant`
CREATE POLICY user_isolation_policy ON users
  USING (slug = current_setting('app.current_user'));

-- The `todos` table has a similar policy restricting rows based on the `current session's tenant`
CREATE POLICY todo_isolation_policy ON todos
  USING (user_slug = current_setting('app.current_user'));

-- Enable RLS on the `users` and `todos` tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE todos ENABLE ROW LEVEL SECURITY;

-- Force RLS enforcement on the `users` and `todos` tables, even for the table owner
ALTER TABLE users FORCE ROW LEVEL SECURITY;
ALTER TABLE todos FORCE ROW LEVEL SECURITY;

-- Seed tables with data:
INSERT INTO users (slug, description)
VALUES
  ('user-1', 'user 1 description'),
  ('user-2', 'user 2 description'),
  ('user-3', 'user 3 description')
ON CONFLICT DO NOTHING;

INSERT INTO todos (title, description, completed, user_slug)
VALUES
  ('My Task #1', 'User 1 Task #1', FALSE, 'user-1'),
  ('My Task #2', 'User 1 Task #2', FALSE, 'user-1'),
  ('My Task #1', 'User 2 Task #1', TRUE, 'user-2'),
  ('My Task #2', 'User 2 Task #2', TRUE, 'user-2');
