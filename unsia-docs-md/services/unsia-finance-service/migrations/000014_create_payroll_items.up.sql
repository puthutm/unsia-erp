-- finance_db: 000014_create_payroll_items.up.sql

CREATE TABLE payroll_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_run_id    UUID NOT NULL REFERENCES payroll_runs(id) ON DELETE CASCADE,
    employee_id       UUID, -- external_ref: hris.employees.id
    gross_amount      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    deduction_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    net_amount        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status            VARCHAR(50) NOT NULL DEFAULT 'draft'
);
