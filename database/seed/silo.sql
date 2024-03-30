-- Seed 1 tenant
INSERT INTO tenants (slug, description)
VALUES
  ('special', 'special tenant using silo db')
ON CONFLICT DO NOTHING;
