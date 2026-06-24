-- Seed data for crm_db

-- Seed Commission Rules
INSERT INTO commission_rules (id, referral_type, amount, calculation_type, is_active, created_at, updated_at)
VALUES 
  ('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'agent', 1000000.00, 'fixed', TRUE, NOW(), NOW()),
  ('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'individual', 500000.00, 'fixed', TRUE, NOW(), NOW()),
  ('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 'public', 150000.00, 'fixed', TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed Active Campaigns
INSERT INTO campaigns (id, code, name, channel, start_date, end_date, status, created_at, updated_at)
VALUES 
  ('d4e5f6a7-b8c9-0d1e-2f3a-4b5c6d7e8f9a', 'CAMP2026', 'PMB Utama 2026', 'Digital Marketing', '2026-01-01', '2026-12-31', 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
