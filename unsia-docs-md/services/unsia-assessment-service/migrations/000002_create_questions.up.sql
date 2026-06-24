-- assessment_db: 000002_create_questions.up.sql

CREATE TABLE questions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_bank_id    UUID NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    question_type       VARCHAR(50) NOT NULL,
    difficulty          VARCHAR(30) NOT NULL DEFAULT 'MEDIUM',
    question_text       TEXT NOT NULL,
    answer_explanation  TEXT,
    status              VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by          UUID, -- external_ref: core.users.id
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
