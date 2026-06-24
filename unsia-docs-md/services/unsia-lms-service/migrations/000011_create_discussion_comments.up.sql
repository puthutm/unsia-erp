-- lms_db: 000011_create_discussion_comments.up.sql

CREATE TABLE discussion_comments (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id      UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id            UUID, -- external_ref: core.users.id
    content            TEXT NOT NULL,
    parent_comment_id  UUID REFERENCES discussion_comments(id) ON DELETE SET NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
