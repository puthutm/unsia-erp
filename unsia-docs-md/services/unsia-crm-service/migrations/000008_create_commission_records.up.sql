-- crm_db: 000008_create_commission_records.up.sql

CREATE TABLE commission_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id             UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    commission_rule_id  UUID REFERENCES commission_rules(id) ON DELETE SET NULL,
    referrer_person_id  UUID, -- external_ref: core.persons.id
    amount              NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status              VARCHAR(50) NOT NULL DEFAULT 'draft',
    sent_to_finance_at  TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
