-- portal_db: 000011_create_portal_contents.up.sql

CREATE TABLE news (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title        VARCHAR(255) NOT NULL,
    content      TEXT NOT NULL,
    author       VARCHAR(100),
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE announcements (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    message     TEXT NOT NULL,
    target_role VARCHAR(100) NOT NULL DEFAULT 'all',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    event_date  DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_news_published ON news(published_at);
CREATE INDEX idx_announcements_role ON announcements(target_role);
CREATE INDEX idx_events_date ON events(event_date);

-- Seed some content
INSERT INTO announcements (title, message, target_role) VALUES
('Selamat Datang', 'Selamat Datang di Portal ERP UNSIA. Gunakan menu navigasi di samping untuk mengakses SIAKAD.', 'all'),
('KRS Online', 'Pengisian KRS Online semester Ganjil 2026/2027 akan ditutup pada 10 Juli 2026.', 'mahasiswa'),
('Batas Penilaian', 'Batas akhir pengisian nilai UAS dosen adalah 15 Juli 2026.', 'dosen');

INSERT INTO news (title, content, author) VALUES
('Lomba Inovasi Nasional', 'ERP UNSIA terpilih sebagai platform terbaik untuk digitalisasi kampus dalam ajang lomba nasional.', 'Humas UNSIA');

INSERT INTO events (title, description, event_date) VALUES
('Mulai Perkuliahan', 'Hari pertama perkuliahan Semester Ganjil 2026/2027.', '2026-09-01');
