-- academic_db: 000003_create_nim_format_configs.up.sql

CREATE TABLE nim_format_configs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             VARCHAR(100) UNIQUE NOT NULL,
    format_template  TEXT,
    token_order      JSONB,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_by       UUID, -- external_ref: core.users.id
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
