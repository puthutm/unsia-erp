-- assessment_db: 000006_create_materials.up.sql

CREATE TABLE materials (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_bank_id  UUID NOT NULL REFERENCES material_banks(id) ON DELETE CASCADE,
    title             VARCHAR(255) NOT NULL,
    material_type     VARCHAR(50) NOT NULL,
    file_url          TEXT,
    content_text      TEXT,
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by        UUID, -- external_ref: core.users.id
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
