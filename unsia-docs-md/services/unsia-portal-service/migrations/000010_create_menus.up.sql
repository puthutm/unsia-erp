-- portal_db: 000010_create_menus.up.sql

CREATE TABLE menus (
    code VARCHAR(100) PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    icon VARCHAR(100),
    parent_code VARCHAR(100) REFERENCES menus(code) ON DELETE SET NULL,
    sort_order INT NOT NULL DEFAULT 0,
    required_permission VARCHAR(100)
);

CREATE TABLE role_menus (
    role_id VARCHAR(100) NOT NULL,
    menu_code VARCHAR(100) REFERENCES menus(code) ON DELETE CASCADE,
    PRIMARY KEY (role_id, menu_code)
);

-- Seed Parent Menus
INSERT INTO menus (code, label, path, icon, parent_code, sort_order) VALUES
('dashboard', 'Dashboard', '/dashboard', 'LayoutDashboard', NULL, 1),
('akademik', 'Akademik', '/akademik', 'GraduationCap', NULL, 2),
('keuangan', 'Keuangan', '/keuangan', 'CircleDollarSign', NULL, 3),
('crm', 'CRM & Marketing', '/crm', 'Users', NULL, 4),
('hris', 'Kepegawaian', '/hris', 'Briefcase', NULL, 5),
('lms', 'LMS (E-Learning)', '/lms', 'BookOpen', NULL, 6),
('cbt', 'Assessment (CBT)', '/cbt', 'FileSpreadsheet', NULL, 7);

-- Seed Submenus
INSERT INTO menus (code, label, path, icon, parent_code, sort_order) VALUES
('mhs_profile', 'Profil Mahasiswa', '/akademik/profil', 'User', 'akademik', 1),
('krs', 'KRS Online', '/akademik/krs', 'FileText', 'akademik', 2),
('khs', 'KHS & Nilai', '/akademik/khs', 'ClipboardList', 'akademik', 3),
('tagihan', 'Tagihan & Pembayaran', '/keuangan/tagihan', 'CreditCard', 'keuangan', 1),
('riwayat_bayar', 'Riwayat Transaksi', '/keuangan/riwayat', 'History', 'keuangan', 2),
('leads', 'Leads Pendaftar', '/crm/leads', 'UserCheck', 'crm', 1),
('agents', 'Agen Referral', '/crm/agents', 'UserPlus', 'crm', 2),
('absensi', 'Presensi Kehadiran', '/hris/absensi', 'Clock', 'hris', 1),
('cuti', 'Pengajuan Cuti', 'calendar', 'hris', 2),
('bkd', 'BKD Dosen', '/hris/bkd', 'FileCheck', 'hris', 3),
('lms_kelas', 'Kelas Aktif', '/lms/kelas', 'Presentation', 'lms', 1),
('lms_tugas', 'Tugas & Penugasan', '/lms/tugas', 'FileSignature', 'lms', 2),
('cbt_jadwal', 'Jadwal Ujian', '/cbt/jadwal', 'CalendarDays', 'cbt', 1),
('cbt_soal', 'Bank Soal', '/cbt/soal', 'Database', 'cbt', 2);

-- Assign Role Menus
-- 1. super_admin has access to everything
INSERT INTO role_menus (role_id, menu_code) 
SELECT 'super_admin', code FROM menus;

-- 2. mahasiswa access
INSERT INTO role_menus (role_id, menu_code) VALUES
('mahasiswa', 'dashboard'),
('mahasiswa', 'akademik'),
('mahasiswa', 'mhs_profile'),
('mahasiswa', 'krs'),
('mahasiswa', 'khs'),
('mahasiswa', 'keuangan'),
('mahasiswa', 'tagihan'),
('mahasiswa', 'riwayat_bayar'),
('mahasiswa', 'lms'),
('mahasiswa', 'lms_kelas'),
('mahasiswa', 'lms_tugas'),
('mahasiswa', 'cbt'),
('mahasiswa', 'cbt_jadwal');

-- 3. dosen access
INSERT INTO role_menus (role_id, menu_code) VALUES
('dosen', 'dashboard'),
('dosen', 'akademik'),
('dosen', 'khs'),
('dosen', 'hris'),
('dosen', 'absensi'),
('dosen', 'cuti'),
('dosen', 'bkd'),
('dosen', 'lms'),
('dosen', 'lms_kelas'),
('dosen', 'cbt'),
('dosen', 'cbt_soal');

-- 4. admin_finance access
INSERT INTO role_menus (role_id, menu_code) VALUES
('admin_finance', 'dashboard'),
('admin_finance', 'keuangan'),
('admin_finance', 'tagihan'),
('admin_finance', 'riwayat_bayar'),
('admin_finance', 'hris'),
('admin_finance', 'absensi'),
('admin_finance', 'cuti');

-- 5. admin_pmb access
INSERT INTO role_menus (role_id, menu_code) VALUES
('admin_pmb', 'dashboard'),
('admin_pmb', 'crm'),
('admin_pmb', 'leads'),
('admin_pmb', 'agents');
