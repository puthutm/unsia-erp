-- lms_db: 000012_create_attendances.up.sql

CREATE TABLE attendances (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id         UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    student_id         UUID NOT NULL, -- external_ref: academic.students.id
    attendance_status  VARCHAR(50) NOT NULL DEFAULT 'present',
    submitted_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_session_student UNIQUE (session_id, student_id)
);

CREATE INDEX idx_lms_attendances_session ON attendances (session_id);
CREATE INDEX idx_lms_attendances_student ON attendances (student_id);
CREATE INDEX idx_lms_attendances_status ON attendances (attendance_status);
