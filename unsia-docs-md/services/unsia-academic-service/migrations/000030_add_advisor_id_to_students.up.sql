-- academic_db: 000030_add_advisor_id_to_students.up.sql

-- Alter students table to add advisor_id and graduation_date columns
ALTER TABLE students ADD COLUMN advisor_id UUID;
ALTER TABLE students ADD COLUMN graduation_date DATE;
CREATE INDEX idx_students_advisor_id ON students (advisor_id);

-- Alter student_advisors table to align with GORM models
ALTER TABLE student_advisors DROP CONSTRAINT IF EXISTS uq_student_period_advisor;

-- Rename lecturer_id to advisor_id
ALTER TABLE student_advisors RENAME COLUMN lecturer_id TO advisor_id;

-- Add academic_year and semester columns
ALTER TABLE student_advisors ADD COLUMN academic_year VARCHAR(50);
ALTER TABLE student_advisors ADD COLUMN semester INT;

-- Drop academic_period_id column (since GORM uses academic_year + semester)
ALTER TABLE student_advisors DROP COLUMN IF EXISTS academic_period_id;

-- Add unique constraint on (student_id, academic_year, semester) to prevent duplicate assignments
ALTER TABLE student_advisors ADD CONSTRAINT uq_student_year_semester_advisor UNIQUE (student_id, academic_year, semester);
