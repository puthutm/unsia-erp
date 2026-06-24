-- lms_db: 000013_create_learning_progress.up.sql

CREATE TABLE learning_progress (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id     UUID UNIQUE NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    progress_percent  NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    last_access_at    TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ
);

CREATE INDEX idx_progress_enrollment ON learning_progress (enrollment_id);
