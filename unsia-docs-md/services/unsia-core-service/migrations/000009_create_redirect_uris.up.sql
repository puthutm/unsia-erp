-- core_db: 000009_create_redirect_uris.up.sql

CREATE TABLE redirect_uris (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    oauth_client_id UUID NOT NULL REFERENCES oauth_clients(id) ON DELETE CASCADE,
    redirect_uri    TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_oauth_client_redirect_uri UNIQUE (oauth_client_id, redirect_uri)
);

CREATE INDEX idx_redirect_uris_oauth_client_id ON redirect_uris (oauth_client_id);
