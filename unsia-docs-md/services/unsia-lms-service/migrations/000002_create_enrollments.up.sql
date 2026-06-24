-- lms_db: 000002_create_enrollments.up.sql

CREATE TABLE enrollments (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lms_class_id       UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    student_id         UUID NOT NULL, -- external_ref: academic.students.id
    enrollment_status  VARCHAR(50) NOT NULL DEFAULT 'active',
    enrolled_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_class_student UNIQUE (lms_class_id, student_id)
);

CREATE INDEX idx_enrollments_class ON enrollments (lms_class_id);
CREATE INDEX idx_enrollments_student ON enrollments (student_id);
CREATE INDEX idx_enrollments_status ON enrollments (enrollment_status);
