-- lms_db: 000010_create_discussions.up.sql

CREATE TABLE discussions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id  UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    created_by  UUID, -- external_ref: core.users.id
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
