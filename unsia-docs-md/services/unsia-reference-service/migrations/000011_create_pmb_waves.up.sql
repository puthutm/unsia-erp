-- reference_db: 000011_create_pmb_waves.up.sql

CREATE TABLE pmb_waves (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year_id            UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    target_entry_period_id      UUID NOT NULL REFERENCES academic_periods(id) ON DELETE CASCADE,
    admission_path_id           UUID REFERENCES admission_paths(id) ON DELETE SET NULL,
    code                        VARCHAR(50) UNIQUE NOT NULL,
    name                        VARCHAR(255) NOT NULL,
    start_date                  DATE,
    end_date                    DATE,
    registration_start_at       TIMESTAMPTZ,
    registration_end_at         TIMESTAMPTZ,
    selection_start_at          TIMESTAMPTZ,
    selection_end_at            TIMESTAMPTZ,
    reregistration_deadline_at  TIMESTAMPTZ,
    status                      VARCHAR(30) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'closed', 'archived')),
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pmb_waves_academic_year_id ON pmb_waves(academic_year_id);
CREATE INDEX idx_pmb_waves_target_entry_period_id ON pmb_waves(target_entry_period_id);
CREATE INDEX idx_pmb_waves_admission_path_id ON pmb_waves(admission_path_id);
CREATE INDEX idx_pmb_waves_status ON pmb_waves(status);
