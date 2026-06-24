-- core_db: 000010_create_service_tokens.up.sql

CREATE TABLE service_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id  UUID NOT NULL REFERENCES applications(id),
    token_hash      TEXT UNIQUE NOT NULL,
    scopes          JSONB,
    expired_at      TIMESTAMPTZ,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_tokens_application_id ON service_tokens (application_id);
CREATE INDEX idx_service_tokens_expired_at ON service_tokens (expired_at);
