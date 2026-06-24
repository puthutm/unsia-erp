-- lms_db: 000004_create_materials.up.sql

CREATE TABLE materials (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id              UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    assessment_material_id  UUID, -- external_ref: assessment.materials.id
    title                   VARCHAR(255) NOT NULL,
    content_type            VARCHAR(100),
    file_url                TEXT,
    published_at            TIMESTAMPTZ
);
