-- portal_db: 000001_create_notifications.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE notifications (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL, -- external_ref: core.users.id
    title          VARCHAR(255) NOT NULL,
    message        TEXT,
    module_source  VARCHAR(100),
    target_url     TEXT,
    sent_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications (user_id);
CREATE INDEX idx_notifications_sent_at ON notifications (sent_at);
