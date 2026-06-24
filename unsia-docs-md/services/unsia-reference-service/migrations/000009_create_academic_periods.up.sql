-- reference_db: 000009_create_academic_periods.up.sql

CREATE TABLE academic_periods (
    id                           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year_id             UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    code                         VARCHAR(50) UNIQUE NOT NULL,
    name                         VARCHAR(255) NOT NULL,
    semester_type                VARCHAR(30) NOT NULL CHECK (semester_type IN ('ganjil', 'genap', 'pendek')),
    start_date                   DATE,
    end_date                     DATE,
    class_start_date             DATE,
    class_end_date               DATE,
    total_meetings               INT DEFAULT 16,
    min_attendance_percentage    NUMERIC(5,2) DEFAULT 75.00,
    status                       VARCHAR(30) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'active', 'closed', 'archived')),
    is_active                    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_academic_year_semester UNIQUE (academic_year_id, semester_type)
);

CREATE INDEX idx_academic_periods_academic_year_id ON academic_periods(academic_year_id);
CREATE INDEX idx_academic_periods_semester_type ON academic_periods(semester_type);
CREATE INDEX idx_academic_periods_status ON academic_periods(status);
