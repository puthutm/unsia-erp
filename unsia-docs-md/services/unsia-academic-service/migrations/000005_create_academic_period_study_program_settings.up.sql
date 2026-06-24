-- academic_db: 000005_create_academic_period_study_program_settings.up.sql

CREATE TABLE academic_period_study_program_settings (
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_period_id         UUID NOT NULL, -- external_ref: ref.academic_periods.id
    study_program_id           UUID NOT NULL, -- external_ref: ref.study_programs.id
    class_start_date           DATE,
    class_end_date             DATE,
    total_meetings             INT DEFAULT 16,
    min_attendance_percentage  NUMERIC(5,2) DEFAULT 75.00,
    pddikti_start_date         DATE,
    pddikti_end_date           DATE,
    is_active                  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_period_study_program UNIQUE (academic_period_id, study_program_id)
);

CREATE INDEX idx_period_sp_settings_period ON academic_period_study_program_settings (academic_period_id);
CREATE INDEX idx_period_sp_settings_sp ON academic_period_study_program_settings (study_program_id);
