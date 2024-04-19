-- Seed 1 tenant
INSERT INTO users (slug)
VALUES
  ('special')
ON CONFLICT DO NOTHING;
