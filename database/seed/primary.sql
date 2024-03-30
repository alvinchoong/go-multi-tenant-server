INSERT INTO tenants (slug, description)
VALUES
  ('normal-abc', 'normal tenant'),
  ('normal-def', 'normal tenant')
ON CONFLICT DO NOTHING;
