-- finance_db: 000009_create_student_clearances.up.sql

CREATE TABLE student_clearances (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL, -- external_ref: academic.students.id
    academic_period_id  UUID, -- external_ref: ref.academic_periods.id
    service_scope       VARCHAR(100) NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'cleared',
    reason              TEXT,
    valid_until         DATE,
    updated_by          UUID, -- external_ref: core.users.id
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
