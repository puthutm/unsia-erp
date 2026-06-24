-- assessment_db: 000007_create_question_sets.up.sql

CREATE TABLE question_sets (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                 VARCHAR(50) UNIQUE NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    randomize_questions  BOOLEAN NOT NULL DEFAULT FALSE,
    randomize_options    BOOLEAN NOT NULL DEFAULT FALSE,
    status               VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
