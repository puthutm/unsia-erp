-- assessment_db: 000003_create_question_versions.up.sql

CREATE TABLE question_versions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id         UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    version_number      INT NOT NULL,
    question_type       VARCHAR(50) NOT NULL,
    difficulty          VARCHAR(30) NOT NULL,
    question_text       TEXT NOT NULL,
    answer_explanation  TEXT,
    options_snapshot    JSONB,
    status              VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_by          UUID, -- external_ref: core.users.id
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_question_version UNIQUE (question_id, version_number)
);

CREATE INDEX idx_question_versions_question ON question_versions (question_id);
