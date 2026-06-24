-- academic_db: 000011_create_class_package_items.up.sql

CREATE TABLE class_package_items (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_package_id      UUID NOT NULL REFERENCES class_packages(id) ON DELETE CASCADE,
    course_id             UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    recommended_class_id  UUID, -- external_ref: academic.classes.id (uuid biasa, no strict FK)
    
    CONSTRAINT uq_package_course UNIQUE (class_package_id, course_id)
);

CREATE INDEX idx_class_package_items_package ON class_package_items (class_package_id);
CREATE INDEX idx_class_package_items_course ON class_package_items (course_id);
