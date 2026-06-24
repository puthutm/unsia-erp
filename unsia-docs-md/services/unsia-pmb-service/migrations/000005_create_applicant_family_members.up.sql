-- pmb_db: 000005_create_applicant_family_members.up.sql

CREATE TABLE applicant_family_members (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id        UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    relation            VARCHAR(50) NOT NULL,
    nik                 VARCHAR(50),
    full_name           VARCHAR(255) NOT NULL,
    education_level_id  UUID, -- external_ref: ref.education_levels.id
    occupation          VARCHAR(100),
    income_range        VARCHAR(100),
    phone               VARCHAR(50),
    dependent_count     INT
);
