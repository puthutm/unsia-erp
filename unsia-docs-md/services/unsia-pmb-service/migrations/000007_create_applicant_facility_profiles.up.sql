-- pmb_db: 000007_create_applicant_facility_profiles.up.sql

CREATE TABLE applicant_facility_profiles (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id         UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    employment_status    VARCHAR(100),
    has_vehicle          BOOLEAN NOT NULL DEFAULT FALSE,
    has_pjj_device       BOOLEAN NOT NULL DEFAULT FALSE,
    internet_access      VARCHAR(100),
    special_need_status  VARCHAR(100),
    CONSTRAINT uq_applicant_facility_profile_applicant UNIQUE (applicant_id)
);
