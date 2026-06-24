-- core_db: 000014_create_audit_logs.up.sql

CREATE TABLE audit_logs (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                     UUID REFERENCES users(id),
    actor_user_id               UUID REFERENCES users(id),
    target_user_id              UUID REFERENCES users(id),
    active_role_id              UUID REFERENCES roles(id),
    impersonation_session_id    UUID REFERENCES impersonation_sessions(id),
    application_id              UUID REFERENCES applications(id),
    module                      VARCHAR(50) NOT NULL,
    action                      VARCHAR(100) NOT NULL,
    entity_name                 VARCHAR(100),
    entity_id                   UUID,
    reason                      TEXT,
    old_value                   JSONB,
    new_value                   JSONB,
    request_id                  VARCHAR(100),
    ip_address                  VARCHAR(50),
    user_agent                  TEXT,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs (actor_user_id);
CREATE INDEX idx_audit_logs_target_user_id ON audit_logs (target_user_id);
CREATE INDEX idx_audit_logs_active_role_id ON audit_logs (active_role_id);
CREATE INDEX idx_audit_logs_impersonation_session_id ON audit_logs (impersonation_session_id);
CREATE INDEX idx_audit_logs_application_id ON audit_logs (application_id);
CREATE INDEX idx_audit_logs_module_entity ON audit_logs (module, entity_name, entity_id);
CREATE INDEX idx_audit_logs_request_id ON audit_logs (request_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);
