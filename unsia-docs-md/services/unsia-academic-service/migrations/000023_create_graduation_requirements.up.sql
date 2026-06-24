-- academic_db: 000023_create_graduation_requirements.up.sql

CREATE TABLE graduation_requirements (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    study_program_id  UUID, -- external_ref: ref.study_programs.id
    degree_level      VARCHAR(50),
    minimum_sks       INT NOT NULL DEFAULT 144,
    minimum_gpa       NUMERIC(3,2) NOT NULL DEFAULT 2.00,
    requirement_json  JSONB,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
