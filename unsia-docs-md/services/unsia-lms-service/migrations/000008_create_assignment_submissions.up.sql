-- lms_db: 000008_create_assignment_submissions.up.sql

CREATE TABLE assignment_submissions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id  UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id     UUID NOT NULL, -- external_ref: academic.students.id
    file_url       TEXT,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    score          NUMERIC(5,2),
    graded_by      UUID, -- external_ref: core.users.id
    graded_at      TIMESTAMPTZ,
    
    CONSTRAINT uq_assignment_student UNIQUE (assignment_id, student_id)
);

CREATE INDEX idx_submissions_assignment ON assignment_submissions (assignment_id);
CREATE INDEX idx_submissions_student ON assignment_submissions (student_id);
