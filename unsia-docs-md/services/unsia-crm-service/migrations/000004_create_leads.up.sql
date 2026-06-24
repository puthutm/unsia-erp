-- crm_db: 000004_create_leads.up.sql

CREATE TABLE leads (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id         UUID NOT NULL, -- external_ref: core.persons.id
    study_program_id  UUID, -- external_ref: ref.study_programs.id
    lead_source_id    UUID, -- external_ref: ref.lead_sources.id
    campaign_id       UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    referral_id       UUID REFERENCES referrals(id) ON DELETE SET NULL,
    lead_number       VARCHAR(50) UNIQUE NOT NULL,
    status            VARCHAR(50) NOT NULL DEFAULT 'new',
    owner_user_id     UUID, -- external_ref: core.users.id
    converted_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_person_id ON leads(person_id);
CREATE INDEX idx_leads_status ON leads(status);
