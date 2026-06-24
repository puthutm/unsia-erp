-- academic_db: 000008_create_courses.up.sql

CREATE TABLE courses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    study_program_id  UUID, -- external_ref: ref.study_programs.id
    course_code       VARCHAR(50) UNIQUE NOT NULL,
    course_name       VARCHAR(255) NOT NULL,
    sks               INT NOT NULL DEFAULT 2,
    course_type       VARCHAR(50), -- e.g. wajib, pilihan, umum
    minimum_grade     NUMERIC(3,2), -- e.g. 2.00 (C)
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_study_program ON courses (study_program_id);
