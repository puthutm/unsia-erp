-- hris_db: 000001_create_work_units.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE work_units (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    parent_unit_id  UUID REFERENCES work_units(id) ON DELETE SET NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);
