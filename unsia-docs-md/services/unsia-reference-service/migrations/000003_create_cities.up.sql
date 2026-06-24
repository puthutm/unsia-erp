-- reference_db: 000003_create_cities.up.sql

CREATE TABLE cities (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    province_id  UUID NOT NULL REFERENCES provinces(id) ON DELETE CASCADE,
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    type         VARCHAR(50),
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_city_code UNIQUE (province_id, code)
);
