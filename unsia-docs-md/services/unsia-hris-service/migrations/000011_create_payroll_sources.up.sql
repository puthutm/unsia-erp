-- hris_db: 000011_create_payroll_sources.up.sql

CREATE TABLE payroll_sources (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id       UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    payroll_period    VARCHAR(50) NOT NULL,
    base_salary       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    allowance_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    deduction_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
