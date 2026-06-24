-- lms_db: 000003_create_sessions.up.sql

CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lms_class_id    UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    session_number  INT NOT NULL,
    title           VARCHAR(255) NOT NULL,
    session_date    DATE,
    start_time      TIME,
    end_time        TIME,
    status          VARCHAR(50) NOT NULL DEFAULT 'draft',
    
    CONSTRAINT uq_class_session UNIQUE (lms_class_id, session_number)
);

CREATE INDEX idx_sessions_class ON sessions (lms_class_id);
CREATE INDEX idx_sessions_date ON sessions (session_date);
CREATE INDEX idx_sessions_status ON sessions (status);
