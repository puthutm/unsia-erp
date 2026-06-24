-- finance_db: 000006_create_scholarships.up.sql

CREATE TABLE scholarships (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id        UUID NOT NULL, -- external_ref: academic.students.id
    scholarship_type  VARCHAR(100),
    amount            NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    approved_by       UUID, -- external_ref: core.users.id
    approved_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
