-- lms_db: 000014_create_grade_syncs.up.sql

CREATE TABLE grade_syncs (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lms_class_id       UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_class_id  UUID NOT NULL, -- external_ref: academic.classes.id
    sync_status        VARCHAR(50) NOT NULL DEFAULT 'pending',
    synced_at          TIMESTAMPTZ,
    payload            JSONB
);
