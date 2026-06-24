-- academic_db: 000020_create_khs.up.sql

CREATE TABLE khs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_period_id  UUID NOT NULL, -- external_ref: ref.academic_periods.id
    ips                 NUMERIC(3,2),
    total_sks           INT,
    file_url            TEXT,
    issued_at           TIMESTAMPTZ,
    
    CONSTRAINT uq_student_period_khs UNIQUE (student_id, academic_period_id)
);

CREATE INDEX idx_khs_student ON khs (student_id);
CREATE INDEX idx_khs_period ON khs (academic_period_id);
