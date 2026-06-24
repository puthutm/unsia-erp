-- academic_db: 000014_create_class_lecturers.up.sql

CREATE TABLE class_lecturers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id     UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    lecturer_id  UUID NOT NULL, -- external_ref: hris.lecturers.id
    role_type    VARCHAR(50) NOT NULL DEFAULT 'teacher', -- coordinator, teacher, assistant

    CONSTRAINT uq_class_lecturer_role UNIQUE (class_id, lecturer_id, role_type)
);

CREATE INDEX idx_class_lecturers_class ON class_lecturers (class_id);
CREATE INDEX idx_class_lecturers_lecturer ON class_lecturers (lecturer_id);
