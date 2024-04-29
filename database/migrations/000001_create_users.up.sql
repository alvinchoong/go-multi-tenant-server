CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
  slug citext PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);
