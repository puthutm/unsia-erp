-- assessment_db: 000008_create_question_set_items.up.sql

CREATE TABLE question_set_items (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_set_id  UUID NOT NULL REFERENCES question_sets(id) ON DELETE CASCADE,
    question_id      UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    score_weight     NUMERIC(5,2) NOT NULL DEFAULT 1.00,
    sort_order       INT NOT NULL DEFAULT 0,
    
    CONSTRAINT uq_set_question UNIQUE (question_set_id, question_id)
);
