-- crm_db: 000001_create_campaigns.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE campaigns (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    channel     VARCHAR(100),
    start_date  DATE,
    end_date    DATE,
    status      VARCHAR(50),
    created_by  UUID, -- external_ref: core.users.id
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
