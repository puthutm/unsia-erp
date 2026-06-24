-- assessment_db: 000001_create_question_banks.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE question_banks (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code          VARCHAR(50) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    module_scope  VARCHAR(100),
    owner_user_id UUID, -- external_ref: core.users.id
    status        VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
