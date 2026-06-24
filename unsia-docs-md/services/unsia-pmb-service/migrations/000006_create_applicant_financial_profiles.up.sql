-- pmb_db: 000006_create_applicant_financial_profiles.up.sql

CREATE TABLE applicant_financial_profiles (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id           UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    personal_income_range  VARCHAR(100),
    bank_name              VARCHAR(100),
    bank_account_name      VARCHAR(255),
    bank_account_number    VARCHAR(100),
    scholarship_interest   BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_applicant_financial_profile_applicant UNIQUE (applicant_id)
);
