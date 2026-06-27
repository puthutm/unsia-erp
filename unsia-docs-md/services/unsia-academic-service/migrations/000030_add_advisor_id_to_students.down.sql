-- academic_db: 000030_add_advisor_id_to_students.down.sql

-- Drop the unique constraint
ALTER TABLE student_advisors DROP CONSTRAINT IF EXISTS uq_student_year_semester_advisor;

-- Add academic_period_id back
ALTER TABLE student_advisors ADD COLUMN academic_period_id UUID;

-- Drop academic_year and semester columns
ALTER TABLE student_advisors DROP COLUMN IF EXISTS academic_year;
ALTER TABLE student_advisors DROP COLUMN IF EXISTS semester;

-- Rename advisor_id back to lecturer_id
ALTER TABLE student_advisors RENAME COLUMN advisor_id TO lecturer_id;

-- Add constraint uq_student_period_advisor back
ALTER TABLE student_advisors ADD CONSTRAINT uq_student_period_advisor UNIQUE (student_id, academic_period_id);

-- Drop advisor_id and graduation_date from students
DROP INDEX IF EXISTS idx_students_advisor_id;
ALTER TABLE students DROP COLUMN IF EXISTS advisor_id;
ALTER TABLE students DROP COLUMN IF EXISTS graduation_date;
