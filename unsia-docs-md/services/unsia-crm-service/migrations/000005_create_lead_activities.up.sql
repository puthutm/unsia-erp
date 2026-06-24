-- crm_db: 000005_create_lead_activities.up.sql

CREATE TABLE lead_activities (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id        UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    user_id        UUID, -- external_ref: core.users.id
    activity_type  VARCHAR(100),
    note           TEXT,
    activity_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
