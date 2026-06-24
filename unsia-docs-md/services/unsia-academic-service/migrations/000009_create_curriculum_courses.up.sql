-- academic_db: 000009_create_curriculum_courses.up.sql

CREATE TABLE curriculum_courses (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    curriculum_id  UUID NOT NULL REFERENCES curriculums(id) ON DELETE CASCADE,
    course_id      UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    semester       INT NOT NULL,
    is_mandatory   BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT uq_curriculum_course UNIQUE (curriculum_id, course_id)
);

CREATE INDEX idx_curr_courses_curriculum ON curriculum_courses (curriculum_id);
CREATE INDEX idx_curr_courses_course ON curriculum_courses (course_id);
