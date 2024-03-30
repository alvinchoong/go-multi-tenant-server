CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE tenants (
    slug citext PRIMARY KEY,
    description TEXT,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);
