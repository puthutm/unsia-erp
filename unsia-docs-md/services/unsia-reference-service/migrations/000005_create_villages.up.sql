-- reference_db: 000005_create_villages.up.sql

CREATE TABLE villages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    district_id  UUID NOT NULL REFERENCES districts(id) ON DELETE CASCADE,
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    postal_code  VARCHAR(20),
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_village_code UNIQUE (district_id, code)
);
