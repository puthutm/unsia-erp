-- lms_db: 000006_create_vicon_links.up.sql

CREATE TABLE vicon_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id  UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    provider    VARCHAR(100),
    join_url    TEXT NOT NULL,
    start_at    TIMESTAMPTZ,
    end_at      TIMESTAMPTZ
);
