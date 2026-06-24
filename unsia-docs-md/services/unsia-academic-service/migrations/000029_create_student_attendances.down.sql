-- Rollback student_attendances table
-- Migration: 000029_create_student_attendances.down.sql

-- Drop foreign key constraints first
ALTER TABLE student_attendances DROP CONSTRAINT IF EXISTS fk_student_attendances_student;
ALTER TABLE student_attendances DROP CONSTRAINT IF EXISTS fk_student_attendances_class;

-- Drop table
DROP TABLE IF EXISTS student_attendances;
