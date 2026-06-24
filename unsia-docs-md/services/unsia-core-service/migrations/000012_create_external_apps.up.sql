-- core_db: 000012_create_external_apps.up.sql

CREATE TABLE external_apps (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    slug            TEXT UNIQUE NOT NULL,
    type            TEXT NOT NULL, -- web, mobile, desktop, webhook, cron
    url             TEXT,
    callback_url    TEXT,
    logo_url        TEXT,
    description    TEXT,
    client_id       TEXT UNIQUE NOT NULL,
    client_secret  TEXT NOT NULL,
    is_active      BOOLEAN NOT NULL DEFAULT true,
    is_internal   BOOLEAN NOT NULL DEFAULT false,
    scopes         JSONB DEFAULT '[]'::jsonb,
    ip_whitelist  JSONB DEFAULT '[]'::jsonb,
    rate_limit    INTEGER DEFAULT 100,
    last_login_at TIMESTAMPTZ,
    expired_at     TIMESTAMPTZ,
    created_by     TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_external_apps_slug ON external_apps (slug);
CREATE INDEX idx_external_apps_client_id ON external_apps (client_id);
CREATE INDEX idx_external_apps_is_active ON external_apps (is_active);
CREATE INDEX idx_external_apps_type ON external_apps (type);
