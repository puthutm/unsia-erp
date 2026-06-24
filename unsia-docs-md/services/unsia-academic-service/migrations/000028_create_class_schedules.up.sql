-- Create class_schedules table
-- Migration: 000028_create_class_schedules.up.sql

CREATE TABLE IF NOT EXISTS class_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 1 AND day_of_week <= 7),
    start_time VARCHAR(5) NOT NULL CHECK (start_time ~ '^[0-2][0-9]:[0-5][0-9]$'),
    end_time VARCHAR(5) NOT NULL CHECK (end_time ~ '^[0-2][0-9]:[0-5][0-9]$'),
    room_id UUID,
    building_id UUID,
    schedule_type VARCHAR(20) DEFAULT 'regular',
    is_online BOOLEAN DEFAULT FALSE,
    meeting_link TEXT,
    capacity_used INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_class_schedules_class_id ON class_schedules(class_id);
CREATE INDEX idx_class_schedules_day_of_week ON class_schedules(day_of_week);
CREATE INDEX idx_class_schedules_room_id ON class_schedules(room_id);
CREATE INDEX idx_class_schedules_building_id ON class_schedules(building_id);

-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    building_id UUID NOT NULL,
    room_code VARCHAR(20) UNIQUE NOT NULL,
    room_name VARCHAR(100) NOT NULL,
    capacity INTEGER DEFAULT 40,
    room_type VARCHAR(20),
    floor INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    equipment JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_rooms_building_id ON rooms(building_id);
CREATE INDEX idx_rooms_room_code ON rooms(room_code);
CREATE INDEX idx_rooms_is_active ON rooms(is_active);

-- Add foreign key constraints
ALTER TABLE class_schedules 
    ADD CONSTRAINT fk_class_schedules_class 
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE;

ALTER TABLE class_schedules 
    ADD CONSTRAINT fk_class_schedules_room 
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL;

ALTER TABLE class_schedules 
    ADD CONSTRAINT fk_class_schedules_building 
    FOREIGN KEY (building_id) REFERENCES ref.gedung(id) ON DELETE SET NULL;

ALTER TABLE rooms 
    ADD CONSTRAINT fk_rooms_building 
    FOREIGN KEY (building_id) REFERENCES ref.gedung(id) ON DELETE RESTRICT;

-- Add comments
COMMENT ON TABLE class_schedules IS 'Jadwal Kelas - Contains class schedule information';
COMMENT ON TABLE rooms IS 'Ruang Kelas - Contains room/venue information';
COMMENT ON COLUMN class_schedules.day_of_week IS 'Day of week: 1=Monday, 7=Sunday';
COMMENT ON COLUMN class_schedules.schedule_type IS 'Type: regular, praktikum, praktik';
