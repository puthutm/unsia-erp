-- crm_db: 000003_create_referrals.up.sql

CREATE TABLE referrals (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referral_type       VARCHAR(50) NOT NULL,
    referrer_person_id  UUID, -- external_ref: core.persons.id
    agent_id            UUID REFERENCES agents(id) ON DELETE SET NULL,
    referral_code       VARCHAR(50) UNIQUE NOT NULL,
    is_valid            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
