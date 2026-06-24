-- reference_db: 000007_create_study_programs.up.sql

CREATE TABLE study_programs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code          VARCHAR(50) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    degree_level  VARCHAR(50),
    faculty_name  VARCHAR(255),
    mode          VARCHAR(50),
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
