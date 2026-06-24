-- Remove foreign key constraint from students table first
ALTER TABLE students DROP CONSTRAINT IF EXISTS fk_students_curriculum;

DROP TABLE IF EXISTS curriculums;
