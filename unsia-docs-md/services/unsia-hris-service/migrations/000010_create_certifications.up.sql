-- hris_db: 000010_create_certifications.up.sql

CREATE TABLE certifications (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id        UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    certification_name VARCHAR(255) NOT NULL,
    issuer             VARCHAR(255) NOT NULL,
    issued_date        DATE,
    expired_date       DATE,
    file_url           TEXT
);
