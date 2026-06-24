-- hris_db: 000005_create_lecturers.up.sql

CREATE TABLE lecturers (
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id                UUID UNIQUE NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    lecturer_status_id         UUID, -- external_ref: ref.lecturer_statuses.id
    functional_position_id     UUID REFERENCES functional_positions(id) ON DELETE SET NULL,
    nidn                       VARCHAR(50) UNIQUE,
    homebase_study_program_id  UUID, -- external_ref: ref.study_programs.id
    certification_status       VARCHAR(50),
    is_active                  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lecturers_employee ON lecturers (employee_id);
CREATE INDEX idx_lecturers_status ON lecturers (lecturer_status_id);
CREATE INDEX idx_lecturers_active ON lecturers (is_active);
