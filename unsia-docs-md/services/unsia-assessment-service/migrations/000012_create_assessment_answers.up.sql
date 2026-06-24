-- assessment_db: 000012_create_assessment_answers.up.sql

CREATE TABLE assessment_answers (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id          UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    question_id         UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id  UUID REFERENCES question_options(id) ON DELETE SET NULL,
    answer_text         TEXT,
    score               NUMERIC(5,2),
    graded_by           UUID, -- external_ref: core.users.id
    graded_at           TIMESTAMPTZ,
    
    CONSTRAINT uq_attempt_question UNIQUE (attempt_id, question_id)
);
