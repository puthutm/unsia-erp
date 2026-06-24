-- core_db: 000025_create_oauth_refresh_tokens.up.sql
-- Kiro spec B-8: Refresh tokens with rotation support

CREATE TABLE oauth_refresh_tokens (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash          TEXT UNIQUE NOT NULL,
    client_id           UUID NOT NULL REFERENCES oauth_clients(id),
    user_id             UUID REFERENCES users(id),  -- nullable for client_credentials (though CC doesn't issue refresh)
    access_token_jti    UUID NOT NULL,
    expires_at          TIMESTAMPTZ NOT NULL,
    used_at             TIMESTAMPTZ,        -- set when rotated (single-use)
    revoked_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oauth_refresh_tokens_client_id ON oauth_refresh_tokens (client_id);
CREATE INDEX idx_oauth_refresh_tokens_user_id ON oauth_refresh_tokens (user_id);
CREATE INDEX idx_oauth_refresh_tokens_access_token_jti ON oauth_refresh_tokens (access_token_jti);
CREATE INDEX idx_oauth_refresh_tokens_expires_at ON oauth_refresh_tokens (expires_at);
