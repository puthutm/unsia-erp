-- core_db: 000001_create_persons.up.sql
-- Persons adalah entitas identitas dasar untuk seluruh ekosistem UNSIA.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE persons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255),
    phone           VARCHAR(50),
    identity_number VARCHAR(50),
    gender          VARCHAR(20),
    birth_place     VARCHAR(100),
    birth_date      DATE,
    religion_id     UUID,           -- external_ref: reference_db.religions.id
    country_id      UUID,           -- external_ref: reference_db.countries.id
    province_id     UUID,           -- external_ref: reference_db.provinces.id
    city_id         UUID,           -- external_ref: reference_db.cities.id
    district_id     UUID,           -- external_ref: reference_db.districts.id
    village_id      UUID,           -- external_ref: reference_db.villages.id
    address         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_persons_email ON persons (email);
CREATE INDEX idx_persons_identity_number ON persons (identity_number);
