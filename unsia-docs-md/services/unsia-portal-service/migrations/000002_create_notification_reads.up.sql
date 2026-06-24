-- portal_db: 000002_create_notification_reads.up.sql

CREATE TABLE notification_reads (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id  UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL, -- external_ref: core.users.id
    read_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_notification_user_read UNIQUE (notification_id, user_id)
);

CREATE INDEX idx_notification_reads_notification ON notification_reads (notification_id);
CREATE INDEX idx_notification_reads_user ON notification_reads (user_id);
