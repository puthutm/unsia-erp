-- academic_db: 000016_create_krs.up.sql

CREATE TABLE krs (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id            UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_period_id    UUID NOT NULL, -- external_ref: ref.academic_periods.id
    status                VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, submitted, approved, rejected
    advisor_id            UUID REFERENCES student_advisors(id) ON DELETE SET NULL,
    is_package            BOOLEAN NOT NULL DEFAULT FALSE,
    finance_clearance_id  UUID, -- external_ref: finance.student_clearances.id
    submitted_at          TIMESTAMPTZ,
    approved_at           TIMESTAMPTZ,
    
    CONSTRAINT uq_student_period_krs UNIQUE (student_id, academic_period_id)
);

CREATE INDEX idx_krs_student ON krs (student_id);
CREATE INDEX idx_krs_period ON krs (academic_period_id);
CREATE INDEX idx_krs_advisor ON krs (advisor_id);
CREATE INDEX idx_krs_status ON krs (status);
