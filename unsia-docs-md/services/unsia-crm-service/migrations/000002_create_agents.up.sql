-- crm_db: 000002_create_agents.up.sql

CREATE TABLE agents (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id          UUID NOT NULL, -- external_ref: core.persons.id
    agent_code         VARCHAR(50) UNIQUE NOT NULL,
    organization_name  VARCHAR(255),
    status             VARCHAR(50) NOT NULL DEFAULT 'active',
    approval_status    VARCHAR(50) NOT NULL DEFAULT 'pending',
    approved_by        UUID, -- external_ref: core.users.id
    approved_at        TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
