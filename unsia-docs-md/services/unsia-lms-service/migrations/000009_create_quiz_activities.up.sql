-- lms_db: 000009_create_quiz_activities.up.sql

CREATE TABLE quiz_activities (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id             UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    assessment_session_id  UUID NOT NULL, -- external_ref: assessment.assessment_sessions.id
    title                  VARCHAR(255) NOT NULL,
    status                 VARCHAR(50) NOT NULL DEFAULT 'active'
);
