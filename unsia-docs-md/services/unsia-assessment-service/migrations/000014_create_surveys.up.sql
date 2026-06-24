-- assessment_db: 000014_create_surveys.up.sql

CREATE TABLE surveys (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title         VARCHAR(255) NOT NULL,
    target_type   VARCHAR(100),
    is_anonymous  BOOLEAN NOT NULL DEFAULT TRUE,
    start_at      TIMESTAMPTZ,
    end_at        TIMESTAMPTZ,
    status        VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by    UUID, -- external_ref: core.users.id
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE survey_questions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id      UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    question_type  VARCHAR(50) NOT NULL,
    question_text  TEXT NOT NULL,
    sort_order     INT NOT NULL DEFAULT 0
);

CREATE TABLE survey_responses (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id           UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    respondent_user_id  UUID, -- external_ref: core.users.id
    submitted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    response_json       JSONB NOT NULL
);
