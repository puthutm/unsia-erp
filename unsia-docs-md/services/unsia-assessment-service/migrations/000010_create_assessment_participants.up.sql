-- assessment_db: 000010_create_assessment_participants.up.sql

CREATE TABLE assessment_participants (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_session_id  UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    participant_type       VARCHAR(50) NOT NULL,
    applicant_id           UUID, -- external_ref: pmb.applicants.id
    student_id             UUID, -- external_ref: academic.students.id
    user_id                UUID, -- external_ref: core.users.id
    status                 VARCHAR(50) NOT NULL DEFAULT 'registered'
);
