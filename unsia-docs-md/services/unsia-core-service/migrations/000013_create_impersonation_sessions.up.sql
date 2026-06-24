-- core_db: 000013_create_impersonation_sessions.up.sql

CREATE TABLE impersonation_sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id       UUID NOT NULL REFERENCES users(id),
    target_user_id      UUID NOT NULL REFERENCES users(id),
    target_role_id      UUID NOT NULL REFERENCES roles(id),
    application_id      UUID REFERENCES applications(id),
    session_id          UUID NOT NULL REFERENCES sessions(id),
    reason              TEXT NOT NULL,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at            TIMESTAMPTZ,
    expired_at          TIMESTAMPTZ NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active'
                        CHECK (status IN ('active', 'ended', 'expired'))
);

CREATE INDEX idx_impersonation_actor ON impersonation_sessions (actor_user_id);
CREATE INDEX idx_impersonation_target ON impersonation_sessions (target_user_id);
CREATE INDEX idx_impersonation_status ON impersonation_sessions (status);
