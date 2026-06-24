-- academic_db: 000010_create_class_packages.up.sql

CREATE TABLE class_packages (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    study_program_id  UUID NOT NULL, -- external_ref: ref.study_programs.id
    curriculum_id     UUID NOT NULL REFERENCES curriculums(id) ON DELETE CASCADE,
    semester          INT NOT NULL,
    package_name      VARCHAR(255) NOT NULL,
    status            VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_class_packages_sp ON class_packages (study_program_id);
CREATE INDEX idx_class_packages_curriculum ON class_packages (curriculum_id);
