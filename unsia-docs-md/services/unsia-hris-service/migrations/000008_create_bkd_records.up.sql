-- hris_db: 000008_create_bkd_records.up.sql

CREATE TABLE bkd_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecturer_id         UUID NOT NULL REFERENCES lecturers(id) ON DELETE CASCADE,
    academic_period_id  UUID, -- external_ref: ref.academic_periods.id
    teaching_load       NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    research_load       NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    service_load        NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    status              VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
