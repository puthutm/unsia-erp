-- academic_db: 000013_create_classes.up.sql

CREATE TABLE classes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_offering_id  UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    class_code          VARCHAR(50) NOT NULL,
    quota               INT NOT NULL DEFAULT 40,
    enrolled_count      INT NOT NULL DEFAULT 0,
    class_status        VARCHAR(50) NOT NULL DEFAULT 'active',
    is_parallel         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_course_offering_class_code UNIQUE (course_offering_id, class_code)
);

CREATE INDEX idx_classes_course_offering ON classes (course_offering_id);
CREATE INDEX idx_classes_code ON classes (class_code);
