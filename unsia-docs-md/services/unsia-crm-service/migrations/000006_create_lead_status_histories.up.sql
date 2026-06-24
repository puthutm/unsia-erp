-- crm_db: 000006_create_lead_status_histories.up.sql

CREATE TABLE lead_status_histories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id     UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    old_status  VARCHAR(50),
    new_status  VARCHAR(50) NOT NULL,
    changed_by  UUID, -- external_ref: core.users.id
    note        TEXT,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
