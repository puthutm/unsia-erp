-- core_db: 000024_create_oauth_access_tokens.up.sql
-- Kiro spec B-8: Access tokens with JTI for revocation tracking

CREATE TABLE oauth_access_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jti         UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    client_id   UUID NOT NULL REFERENCES oauth_clients(id),
    user_id     UUID REFERENCES users(id),  -- nullable for client_credentials flow
    scope       VARCHAR(500),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oauth_access_tokens_client_id ON oauth_access_tokens (client_id);
CREATE INDEX idx_oauth_access_tokens_user_id ON oauth_access_tokens (user_id);
CREATE INDEX idx_oauth_access_tokens_expires_at ON oauth_access_tokens (expires_at);
CREATE INDEX idx_oauth_access_tokens_revoked_at ON oauth_access_tokens (revoked_at);
