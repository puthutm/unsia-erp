-- academic_db: 000002_create_student_advisors.up.sql

CREATE TABLE student_advisors (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    lecturer_id         UUID NOT NULL, -- external_ref: hris.lecturers.id
    academic_period_id  UUID, -- external_ref: ref.academic_periods.id
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_student_period_advisor UNIQUE (student_id, academic_period_id)
);

CREATE INDEX idx_student_advisors_student ON student_advisors (student_id);
CREATE INDEX idx_student_advisors_lecturer ON student_advisors (lecturer_id);
CREATE INDEX idx_student_advisors_period ON student_advisors (academic_period_id);
