-- Create student_attendances table
-- Migration: 000029_create_student_attendances.up.sql

CREATE TABLE IF NOT EXISTS student_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    class_id UUID NOT NULL,
    session_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'excused', 'sick')),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_student_attendances_student_id ON student_attendances(student_id);
CREATE INDEX idx_student_attendances_class_id ON student_attendances(class_id);
CREATE INDEX idx_student_attendances_session_date ON student_attendances(session_date);
CREATE UNIQUE INDEX idx_student_attendances_unique ON student_attendances(student_id, class_id, session_date);

-- Add foreign key constraints
ALTER TABLE student_attendances 
    ADD CONSTRAINT fk_student_attendances_student 
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE;

ALTER TABLE student_attendances 
    ADD CONSTRAINT fk_student_attendances_class 
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE;

-- Add comments
COMMENT ON TABLE student_attendances IS 'Absensi Mahasiswa - Tracks student attendance per class session';
COMMENT ON COLUMN student_attendances.status IS 'Status: present, absent, excused, sick';
