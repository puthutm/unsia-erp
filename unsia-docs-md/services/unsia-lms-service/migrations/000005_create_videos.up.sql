-- lms_db: 000005_create_videos.up.sql

CREATE TABLE videos (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id        UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    title             VARCHAR(255) NOT NULL,
    video_url         TEXT NOT NULL,
    duration_minutes  INT
);
