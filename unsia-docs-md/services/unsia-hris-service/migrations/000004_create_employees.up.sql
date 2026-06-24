-- hris_db: 000004_create_employees.up.sql

CREATE TABLE employees (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id          UUID NOT NULL, -- external_ref: core.persons.id
    employee_type_id   UUID, -- external_ref: ref.employee_types.id
    work_unit_id       UUID REFERENCES work_units(id) ON DELETE SET NULL,
    position_id        UUID REFERENCES positions(id) ON DELETE SET NULL,
    nip                VARCHAR(100) UNIQUE,
    employment_status  VARCHAR(50) NOT NULL DEFAULT 'contract',
    join_date          DATE,
    end_date           DATE,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_person ON employees (person_id);
CREATE INDEX idx_employees_status ON employees (employment_status);
CREATE INDEX idx_employees_active ON employees (is_active);
