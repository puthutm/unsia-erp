-- academic_db: 000015_create_class_schedules.up.sql

CREATE TABLE class_schedules (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id      UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    day           VARCHAR(20) NOT NULL, -- e.g. Senin, Selasa
    start_time    TIME NOT NULL,
    end_time      TIME NOT NULL,
    room_or_link  TEXT,
    session_type  VARCHAR(50) NOT NULL DEFAULT 'online' -- online, offline, hybrid
);

CREATE INDEX idx_class_schedules_class ON class_schedules (class_id);
