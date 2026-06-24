-- pmb_db: 000004_create_applicant_education_backgrounds.up.sql

CREATE TABLE applicant_education_backgrounds (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id           UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    education_level_id     UUID, -- external_ref: ref.education_levels.id
    institution_name       VARCHAR(255),
    npsn_or_pt_code        VARCHAR(50),
    nisn_or_previous_nim   VARCHAR(50),
    graduation_year        INT,
    average_score          NUMERIC(5,2)
);
