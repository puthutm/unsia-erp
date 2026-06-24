-- core_db: 000023_create_oauth_authorization_codes.up.sql
-- Kiro spec B-8: Authorization codes for OAuth 2.0 Authorization Code Flow + PKCE

CREATE TABLE oauth_authorization_codes (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_hash               TEXT NOT NULL,
    client_id               UUID NOT NULL REFERENCES oauth_clients(id),
    user_id                 UUID NOT NULL REFERENCES users(id),
    redirect_uri            TEXT NOT NULL,
    scope                   VARCHAR(500),
    code_challenge          VARCHAR(200) NOT NULL,
    code_challenge_method   VARCHAR(10) NOT NULL DEFAULT 'S256'
                            CHECK (code_challenge_method = 'S256'),
    state                   VARCHAR(500),
    expires_at              TIMESTAMPTZ NOT NULL,
    used_at                 TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_oauth_auth_codes_code_hash ON oauth_authorization_codes (code_hash);
CREATE INDEX idx_oauth_auth_codes_client_id ON oauth_authorization_codes (client_id);
CREATE INDEX idx_oauth_auth_codes_user_id ON oauth_authorization_codes (user_id);
CREATE INDEX idx_oauth_auth_codes_expires_at ON oauth_authorization_codes (expires_at);
