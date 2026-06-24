-- finance_db: 000001_create_invoices.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE invoices (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number      VARCHAR(100) UNIQUE NOT NULL,
    target_type         VARCHAR(50) NOT NULL,
    applicant_id        UUID, -- external_ref: pmb.applicants.id
    student_id          UUID, -- external_ref: academic.students.id
    academic_period_id  UUID, -- external_ref: ref.academic_periods.id
    total_amount        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    paid_amount         NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status              VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    due_date            DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
