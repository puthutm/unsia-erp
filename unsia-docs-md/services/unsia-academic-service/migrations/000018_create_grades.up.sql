-- academic_db: 000018_create_grades.up.sql

CREATE TABLE grades (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    krs_item_id    UUID NOT NULL REFERENCES krs_items(id) ON DELETE CASCADE,
    numeric_grade  NUMERIC(5,2),
    letter_grade   VARCHAR(5),
    grade_point    NUMERIC(3,2),
    source         VARCHAR(50), -- e.g. lms, manual
    submitted_at   TIMESTAMPTZ,
    submitted_by   UUID, -- external_ref: core.users.id
    
    CONSTRAINT uq_krs_item_grade UNIQUE (krs_item_id)
);

CREATE INDEX idx_grades_krs_item ON grades (krs_item_id);
CREATE INDEX idx_grades_submitted_by ON grades (submitted_by);
