-- reference_db: 000004_create_districts.up.sql

CREATE TABLE districts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id      UUID NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_district_code UNIQUE (city_id, code)
);
