-- lms_db: 000001_create_classes.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE classes (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_class_id  UUID UNIQUE NOT NULL, -- external_ref: academic.classes.id
    lecturer_id        UUID, -- external_ref: hris.lecturers.id
    status             VARCHAR(50) NOT NULL DEFAULT 'active',
    synced_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_classes_lecturer ON classes (lecturer_id);
