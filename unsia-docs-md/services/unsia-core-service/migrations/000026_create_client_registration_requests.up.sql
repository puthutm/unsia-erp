-- core_db: 000026_create_client_registration_requests.up.sql
-- Kiro spec B-8: Audit trail untuk semua permintaan registrasi OAuth client

CREATE TABLE client_registration_requests (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    oauth_client_id         UUID REFERENCES oauth_clients(id),  -- diisi setelah approved
    owner_name              VARCHAR(255) NOT NULL,
    owner_email             VARCHAR(255) NOT NULL,
    owner_organization      VARCHAR(255) NOT NULL,
    requested_scopes        JSONB,
    requested_grant_types   JSONB,
    requested_redirect_uris JSONB,
    status                  VARCHAR(20) NOT NULL DEFAULT 'PENDING'
                            CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    admin_notes             TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at             TIMESTAMPTZ,
    reviewed_by             UUID
);

CREATE INDEX idx_client_reg_requests_status ON client_registration_requests (status);
CREATE INDEX idx_client_reg_requests_owner_email ON client_registration_requests (owner_email);
CREATE INDEX idx_client_reg_requests_oauth_client_id ON client_registration_requests (oauth_client_id);
