-- create table
CREATE TABLE users (
    id citext,
    tenant_slug citext NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, tenant_slug),
    CONSTRAINT fk_slug FOREIGN KEY (tenant_slug) REFERENCES tenants(slug)
);

-- enable row level security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
-- subject table owner to row security
ALTER TABLE users FORCE ROW LEVEL SECURITY;

-- row level security policy
CREATE POLICY user_isolation_policy ON users
  USING (tenant_slug = current_setting('app.current_tenant'));
