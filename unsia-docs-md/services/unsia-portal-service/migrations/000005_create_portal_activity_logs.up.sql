-- portal_db: 000005_create_portal_activity_logs.up.sql

CREATE TABLE portal_activity_logs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL, -- external_ref: core.users.id
    activity_type  VARCHAR(100) NOT NULL,
    module_target  VARCHAR(100),
    description    TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user ON portal_activity_logs (user_id);
