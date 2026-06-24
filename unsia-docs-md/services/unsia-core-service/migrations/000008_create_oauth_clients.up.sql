-- core_db: 000008_create_oauth_clients.up.sql

CREATE TABLE oauth_clients (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id      UUID NOT NULL REFERENCES applications(id),
    client_id           VARCHAR(255) UNIQUE NOT NULL,
    client_secret_hash  TEXT,
    client_name         VARCHAR(100) NOT NULL,
    client_type         VARCHAR(20) NOT NULL DEFAULT 'confidential'
                        CHECK (client_type IN ('confidential', 'public')),
    grant_types         JSONB NOT NULL DEFAULT '["authorization_code"]',
    allowed_scopes      JSONB NOT NULL DEFAULT '[]',
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING'
                        CHECK (status IN ('PENDING', 'ACTIVE', 'SUSPENDED', 'REVOKED')),
    owner_name          VARCHAR(255),
    owner_email         VARCHAR(255),
    owner_organization  VARCHAR(255),
    approved_at         TIMESTAMPTZ,
    approved_by         UUID,
    suspended_at        TIMESTAMPTZ,
    revoked_at          TIMESTAMPTZ,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oauth_clients_application_id ON oauth_clients (application_id);
CREATE INDEX idx_oauth_clients_client_id ON oauth_clients (client_id);
CREATE INDEX idx_oauth_clients_status ON oauth_clients (status);
CREATE INDEX idx_oauth_clients_owner_email ON oauth_clients (owner_email);
