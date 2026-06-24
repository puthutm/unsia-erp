-- academic_db: 000024_create_yudisium_records.up.sql

CREATE TABLE yudisium_records (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id         UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    yudisium_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    graduation_status  VARCHAR(50) NOT NULL DEFAULT 'lulus',
    final_gpa          NUMERIC(3,2),
    transcript_number  VARCHAR(100),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_student_yudisium UNIQUE (student_id)
);

CREATE INDEX idx_yudisium_student ON yudisium_records (student_id);
