-- finance_db: 000013_create_payroll_runs.up.sql

CREATE TABLE payroll_runs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_period VARCHAR(50) NOT NULL,
    run_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    total_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status        VARCHAR(50) NOT NULL DEFAULT 'draft',
    approved_by   UUID, -- external_ref: core.users.id
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
