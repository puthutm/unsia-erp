-- portal_db: 000003_create_user_preferences.up.sql

CREATE TABLE user_preferences (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL, -- external_ref: core.users.id
    preference_key    VARCHAR(100) NOT NULL,
    preference_value  JSONB NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_user_preference_key UNIQUE (user_id, preference_key)
);

CREATE INDEX idx_user_preferences_user ON user_preferences (user_id);
