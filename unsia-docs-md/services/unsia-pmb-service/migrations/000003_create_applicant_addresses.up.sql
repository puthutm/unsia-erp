-- pmb_db: 000003_create_applicant_addresses.up.sql

CREATE TABLE applicant_addresses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id    UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    address_type    VARCHAR(50),
    street          TEXT,
    province_id     UUID, -- external_ref: ref.provinces.id
    city_id         UUID, -- external_ref: ref.cities.id
    district_id     UUID, -- external_ref: ref.districts.id
    village_id      UUID, -- external_ref: ref.villages.id
    postal_code     VARCHAR(20),
    is_same_as_ktp  BOOLEAN NOT NULL DEFAULT FALSE
);
