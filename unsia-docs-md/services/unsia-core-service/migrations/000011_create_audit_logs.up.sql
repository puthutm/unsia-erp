-- core_db: 000011_create_audit_logs.up.sql

CREATE TABLE audit_logs (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                    UUID REFERENCES users(id),
    actor_user_id               UUID REFERENCES users(id),
    target_user_id              UUID REFERENCES users(id),
    active_role_id              UUID REFERENCES roles(id),
    impersonation_session_id    UUID REFERENCES impersonation_sessions(id),
    application_id             UUID REFERENCES applications(id),
    module                     TEXT NOT NULL,
    action                     TEXT NOT NULL,
    entity_name                TEXT,
    entity_id                  UUID,
    reason                     TEXT,
    old_value                  JSONB,
    new_value                  JSONB,
    request_id                 TEXT,
    ip_address                 INET,
    user_agent                 TEXT,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs (actor_user_id);
CREATE INDEX idx_audit_logs_target_user_id ON audit_logs (target_user_id);
CREATE INDEX idx_audit_logs_module_action ON audit_logs (module, action);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_name, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at DESC);
CREATE INDEX idx_audit_logs_application_id ON audit_logs (application_id);
