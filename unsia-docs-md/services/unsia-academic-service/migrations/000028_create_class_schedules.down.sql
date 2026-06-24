-- Rollback class_schedules and rooms tables
-- Migration: 000028_create_class_schedules.down.sql

-- Drop foreign key constraints first
ALTER TABLE class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_class;
ALTER TABLE class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_room;
ALTER TABLE class_schedules DROP CONSTRAINT IF EXISTS fk_class_schedules_building;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS fk_rooms_building;

-- Drop tables
DROP TABLE IF EXISTS class_schedules;
DROP TABLE IF EXISTS rooms;
