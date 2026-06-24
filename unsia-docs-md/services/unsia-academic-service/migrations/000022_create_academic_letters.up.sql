-- academic_db: 000022_create_academic_letters.up.sql

CREATE TABLE academic_letters (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id    UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    letter_type   VARCHAR(100) NOT NULL, -- e.g. keterangan_aktif, cuti
    status        VARCHAR(50) NOT NULL DEFAULT 'requested', -- requested, issued, rejected
    file_url      TEXT,
    requested_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    issued_at     TIMESTAMPTZ
);

CREATE INDEX idx_academic_letters_student ON academic_letters (student_id);
CREATE INDEX idx_academic_letters_status ON academic_letters (status);
