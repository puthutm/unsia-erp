-- academic_db: 000007_create_curriculums.up.sql

CREATE TABLE curriculums (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    study_program_id            UUID NOT NULL, -- external_ref: ref.study_programs.id
    code                        VARCHAR(100) UNIQUE NOT NULL,
    name                        VARCHAR(255) NOT NULL,
    curriculum_year             INT NOT NULL,
    status                      VARCHAR(50) NOT NULL DEFAULT 'draft',
    effective_start_period_id   UUID, -- external_ref: ref.academic_periods.id
    effective_end_period_id     UUID, -- external_ref: ref.academic_periods.id
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    is_default_for_new_student  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_sp_curriculum_year UNIQUE (study_program_id, curriculum_year)
);

CREATE INDEX idx_curriculums_study_program ON curriculums (study_program_id);
CREATE INDEX idx_curriculums_year ON curriculums (curriculum_year);
CREATE INDEX idx_curriculums_status ON curriculums (status);

-- Add foreign key constraint to students table
ALTER TABLE students 
    ADD CONSTRAINT fk_students_curriculum 
    FOREIGN KEY (curriculum_id) 
    REFERENCES curriculums(id) 
    ON DELETE SET NULL;
