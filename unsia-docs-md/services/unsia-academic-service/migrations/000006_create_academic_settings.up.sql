-- academic_db: 000006_create_academic_settings.up.sql

CREATE TABLE academic_settings (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key    VARCHAR(100) UNIQUE NOT NULL,
    setting_value  JSONB,
    updated_by     UUID, -- external_ref: core.users.id
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
