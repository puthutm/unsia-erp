-- core_db: 000012_create_active_role_sessions.up.sql

CREATE TABLE active_role_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    role_id         UUID NOT NULL REFERENCES roles(id),
    session_id      UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    application_id  UUID REFERENCES applications(id),
    activated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_session_application UNIQUE (session_id, application_id)
);

CREATE INDEX idx_active_role_sessions_user_id ON active_role_sessions (user_id);
CREATE INDEX idx_active_role_sessions_role_id ON active_role_sessions (role_id);
CREATE INDEX idx_active_role_sessions_session_id ON active_role_sessions (session_id);
CREATE INDEX idx_active_role_sessions_application_id ON active_role_sessions (application_id);
