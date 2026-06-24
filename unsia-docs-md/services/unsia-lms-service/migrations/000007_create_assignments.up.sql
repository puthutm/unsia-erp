-- lms_db: 000007_create_assignments.up.sql

CREATE TABLE assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    instruction  TEXT,
    due_at       TIMESTAMPTZ,
    status       VARCHAR(50) NOT NULL DEFAULT 'active'
);
