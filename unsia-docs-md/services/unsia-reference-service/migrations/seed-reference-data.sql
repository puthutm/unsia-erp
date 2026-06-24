-- Seeding initial master data for reference_db
-- Run this after running migrations on reference_db

-- 1. Seeding Religions
INSERT INTO religions (id, name) VALUES
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Islam'),
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Kristen Protestan'),
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Katolik'),
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'Hindu'),
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a05', 'Buddha'),
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a06', 'Khonghucu')
ON CONFLICT (id) DO NOTHING;

-- 2. Seeding Countries
INSERT INTO countries (id, code, name) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'ID', 'Indonesia')
ON CONFLICT (id) DO NOTHING;

-- 3. Seeding Provinces (Indonesia)
INSERT INTO provinces (id, country_id, name) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'DKI Jakarta'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Jawa Barat'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Banten')
ON CONFLICT (id) DO NOTHING;

-- 4. Seeding Cities
INSERT INTO cities (id, province_id, name) VALUES
('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Jakarta Selatan'),
('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Bandung'),
('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Tangerang')
ON CONFLICT (id) DO NOTHING;

-- 5. Seeding Study Programs
INSERT INTO study_programs (id, code, name, degree, status, created_at, updated_at) VALUES
('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'IF', 'Informatika', 'S1', 'ACTIVE', NOW(), NOW()),
('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'SI', 'Sistem Informasi', 'S1', 'ACTIVE', NOW(), NOW()),
('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'MN', 'Manajemen', 'S1', 'ACTIVE', NOW(), NOW()),
('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'IK', 'Ilmu Komunikasi', 'S1', 'ACTIVE', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 6. Seeding Status Codes (Managed Business State Machines)
INSERT INTO status_codes (id, module, code, name, description) VALUES
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'pmb', 'DRAFT', 'Draft', 'Pendaftaran baru dibuat oleh pemohon'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'pmb', 'SUBMITTED', 'Submitted', 'Pendaftaran telah disubmit untuk verifikasi berkas'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'pmb', 'APPROVED', 'Approved', 'Berkas pendaftaran dinyatakan lolos seleksi'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'pmb', 'REJECTED', 'Rejected', 'Berkas ditolak, pendaftar harus memperbaiki data'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a05', 'pmb', 'HANDED_OVER', 'Handed Over', 'Data mahasiswa sukses diserahkan ke Akademik'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a06', 'finance', 'UNPAID', 'Unpaid', 'Invoice belum terbayarkan'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a07', 'finance', 'PAID', 'Paid', 'Invoice telah lunas terbayar'),
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a08', 'finance', 'EXPIRED', 'Expired', 'Masa tenggang pembayaran invoice habis')
ON CONFLICT (id) DO NOTHING;

-- 7. Seeding Payment Components
INSERT INTO payment_components (id, code, name, default_amount, is_active, created_at, updated_at) VALUES
('71eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'pmb_reg', 'Uang Pendaftaran PMB', 250000.00, TRUE, NOW(), NOW()),
('71eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'development_fee', 'Uang Gedung / Pengembangan', 2000000.00, TRUE, NOW(), NOW()),
('71eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'tuition_sem', 'Uang Kuliah Semester (BPP)', 3000000.00, TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 8. Seeding Payment Methods
INSERT INTO payment_methods (id, code, name, provider, is_active, created_at, updated_at) VALUES
('81eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'va_mandiri', 'Mandiri Virtual Account', 'Midtrans', TRUE, NOW(), NOW()),
('81eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'va_bni', 'BNI Virtual Account', 'Midtrans', TRUE, NOW(), NOW()),
('81eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'credit_card', 'Kartu Kredit', 'Midtrans', TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 9. Seeding Document Types
INSERT INTO document_types (id, code, name, is_mandatory, is_active, created_at, updated_at) VALUES
('91eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'ijazah_sma', 'Ijazah SMA / Sederajat', TRUE, TRUE, NOW(), NOW()),
('91eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ktp', 'Kartu Tanda Penduduk (KTP)', TRUE, TRUE, NOW(), NOW()),
('91eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'kk', 'Kartu Keluarga (KK)', FALSE, TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
