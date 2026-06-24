-- assessment_db: 000009_create_assessment_sessions.up.sql

CREATE TABLE assessment_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_set_id   UUID REFERENCES question_sets(id) ON DELETE SET NULL,
    session_type      VARCHAR(50) NOT NULL,
    context_module    VARCHAR(100),
    context_id        UUID,
    title             VARCHAR(255) NOT NULL,
    start_at          TIMESTAMPTZ,
    end_at            TIMESTAMPTZ,
    duration_minutes  INT,
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    passing_grade     NUMERIC(5,2),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
