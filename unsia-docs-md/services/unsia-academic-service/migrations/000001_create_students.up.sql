-- academic_db: 000001_create_students.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE students (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id               UUID NOT NULL, -- external_ref: core.persons.id
    user_id                 UUID, -- external_ref: core.users.id
    applicant_id            UUID UNIQUE, -- external_ref: pmb.applicants.id
    study_program_id        UUID NOT NULL, -- external_ref: ref.study_programs.id
    nim                     VARCHAR(50) UNIQUE NOT NULL,
    student_status          VARCHAR(50) NOT NULL DEFAULT 'active',
    entry_academic_year_id  UUID, -- external_ref: ref.academic_years.id
    entry_period_id         UUID, -- external_ref: ref.academic_periods.id
    curriculum_id           UUID, 
    current_semester        INT NOT NULL DEFAULT 1,
    active_date             DATE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_students_person ON students (person_id);
CREATE INDEX idx_students_user ON students (user_id);
CREATE INDEX idx_students_study_program ON students (study_program_id);
CREATE INDEX idx_students_status ON students (student_status);
CREATE INDEX idx_students_entry_period ON students (entry_period_id);
