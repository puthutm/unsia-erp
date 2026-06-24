-- assessment_db: 000013_create_assessment_scores.up.sql

CREATE TABLE assessment_scores (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id          UUID UNIQUE NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    total_score         NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    result_status       VARCHAR(50),
    published_at        TIMESTAMPTZ,
    sent_to_context_at  TIMESTAMPTZ
);
