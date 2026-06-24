-- pmb_db: 000001_create_applicants.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE applicants (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id               UUID NOT NULL, -- external_ref: core.persons.id
    user_id                 UUID, -- external_ref: core.users.id
    crm_lead_id             UUID, -- external_ref: crm.leads.id
    study_program_id        UUID, -- external_ref: ref.study_programs.id
    pmb_wave_id             UUID, -- external_ref: ref.pmb_waves.id
    admission_path_id       UUID, -- external_ref: ref.admission_paths.id
    target_entry_period_id  UUID, -- external_ref: ref.academic_periods.id
    registration_number     VARCHAR(100) UNIQUE NOT NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'verified', 'accepted', 'reregistration_completed', 'ready_for_academic')),
    submitted_at            TIMESTAMPTZ,
    accepted_at             TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_applicants_person_id ON applicants (person_id);
CREATE INDEX idx_applicants_user_id ON applicants (user_id);
CREATE INDEX idx_applicants_crm_lead_id ON applicants (crm_lead_id);
CREATE INDEX idx_applicants_study_program_id ON applicants (study_program_id);
CREATE INDEX idx_applicants_pmb_wave_id ON applicants (pmb_wave_id);
CREATE INDEX idx_applicants_admission_path_id ON applicants (admission_path_id);
CREATE INDEX idx_applicants_target_entry_period_id ON applicants (target_entry_period_id);
CREATE INDEX idx_applicants_status ON applicants (status);
