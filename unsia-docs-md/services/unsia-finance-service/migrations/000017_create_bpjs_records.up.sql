-- finance_db: 000017_create_bpjs_records.up.sql

CREATE TABLE bpjs_records (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id  UUID, -- external_ref: hris.employees.id
    amount       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    period       VARCHAR(50) NOT NULL,
    status       VARCHAR(50) NOT NULL DEFAULT 'unpaid'
);
