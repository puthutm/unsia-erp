-- ==========================================
-- SUPER ADMIN SEEDER for UNSIA ERP
-- ==========================================
-- Run this after running migrations on core_db
-- 
-- Default login credentials:
-- Username: admin
-- Password: password123
--
-- To execute, run:
-- psql -h localhost -U postgres -d core_db -f seed-super-admin.sql

-- 1. Insert Persons (Super Admin)
INSERT INTO persons (id, name, email, phone, created_at, updated_at)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Super Administrator',
    'admin@unsia.ac.id',
    '081234567890',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO persons (id, name, email, phone, created_at, updated_at)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
    'PMB Admissions Officer',
    'pmb.admin@unsia.ac.id',
    '081234567891',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 2. Insert Users (Password: password123)
-- Bcrypt hash: $2a$10$8K.g32uE9c3nJ/U5uU5jO.HwD3Ym6lE29X0E7oYfLqS.m3o3B6B9e
INSERT INTO users (id, person_id, username, password_hash, status, created_at, updated_at)
VALUES (
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'admin',
    '$2a$10$8K.g32uE9c3nJ/U5uU5jO.HwD3Ym6lE29X0E7oYfLqS.m3o3B6B9e',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, person_id, username, password_hash, status, created_at, updated_at)
VALUES (
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12',
    'admin_pmb',
    '$2a$10$8K.g32uE9c3nJ/U5uU5jO.HwD3Ym6lE29X0E7oYfLqS.m3o3B6B9e',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 3. Insert Roles
INSERT INTO roles (id, code, name, scope_type, created_at, updated_at)
VALUES (
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'super_admin',
    'Super Admin',
    'global',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO roles (id, code, name, scope_type, created_at, updated_at)
VALUES (
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'admin_pmb',
    'Admin PMB',
    'global',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 4. Insert Permissions
INSERT INTO permissions (id, code, name, module, created_at, updated_at)
VALUES (
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
    '*',
    'All Permissions Access',
    'core',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO permissions (id, code, name, module, created_at, updated_at)
VALUES (
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
    'pmb.applicant.verify_document',
    'Verify Applicant Document',
    'pmb',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 5. Insert Role Permissions
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
VALUES (
    'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- super_admin
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', -- *
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO role_permissions (id, role_id, permission_id, created_at)
VALUES (
    'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02',
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- admin_pmb
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', -- pmb.applicant.verify_document
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 6. Insert User Roles
INSERT INTO user_roles (id, user_id, role_id, study_program_id, created_at)
VALUES (
    'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380c01',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- admin user
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- super_admin role
    NULL,
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (id, user_id, role_id, study_program_id, created_at)
VALUES (
    'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380c02',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', -- admin_pmb user
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- admin_pmb role
    NULL,
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 7. Insert Applications Registry
INSERT INTO applications (id, application_code, name, url, enabled, created_at, updated_at)
VALUES (
    '9b0faa18-de61-43a0-ac9b-34f7e80fe7a0',
    'PORTAL',
    'UNSIA Portal',
    'http://localhost:3000',
    TRUE,
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
