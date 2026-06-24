-- academic_db: 000021_create_transcripts.up.sql

CREATE TABLE transcripts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id  UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    ipk         NUMERIC(3,2),
    total_sks   INT,
    file_url    TEXT,
    issued_at   TIMESTAMPTZ,
    
    CONSTRAINT uq_student_transcript UNIQUE (student_id)
);

CREATE INDEX idx_transcripts_student ON transcripts (student_id);
