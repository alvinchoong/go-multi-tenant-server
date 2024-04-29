CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- create table
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

-- enable row level security
ALTER TABLE todos ENABLE ROW LEVEL SECURITY;
-- subject table owner to row security
ALTER TABLE todos FORCE ROW LEVEL SECURITY;

-- row level security policy
CREATE POLICY todo_isolation_policy ON todos
  USING (user_slug = current_setting('app.current_user'));
