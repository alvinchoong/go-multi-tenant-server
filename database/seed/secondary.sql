INSERT INTO tenants (slug, description)
VALUES
  ('special-abc', 'special tenant'),
  ('special-def', 'special tenant')
ON CONFLICT DO NOTHING;
