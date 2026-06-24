-- assessment_db: 000011_create_assessment_attempts.up.sql

CREATE TABLE assessment_attempts (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_session_id  UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    participant_id         UUID NOT NULL REFERENCES assessment_participants(id) ON DELETE CASCADE,
    attempt_number         INT NOT NULL DEFAULT 1,
    idempotency_key        VARCHAR(500) UNIQUE,
    started_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitted_at           TIMESTAMPTZ,
    status                 VARCHAR(50) NOT NULL DEFAULT 'started',
    total_score            NUMERIC(5,2),
    
    CONSTRAINT uq_session_participant_attempt UNIQUE (assessment_session_id, participant_id, attempt_number)
);
