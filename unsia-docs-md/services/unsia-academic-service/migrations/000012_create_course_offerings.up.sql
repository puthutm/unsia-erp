-- academic_db: 000012_create_course_offerings.up.sql

CREATE TABLE course_offerings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id           UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    academic_period_id  UUID NOT NULL, -- external_ref: ref.academic_periods.id
    curriculum_id       UUID REFERENCES curriculums(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'active',
    opened_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_course_period_curriculum UNIQUE (course_id, academic_period_id, curriculum_id)
);

CREATE INDEX idx_course_offerings_course ON course_offerings (course_id);
CREATE INDEX idx_course_offerings_period ON course_offerings (academic_period_id);
CREATE INDEX idx_course_offerings_curriculum ON course_offerings (curriculum_id);
CREATE INDEX idx_course_offerings_status ON course_offerings (status);
