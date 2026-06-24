-- academic_db: 000025_create_alumni.up.sql

CREATE TABLE alumni (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id       UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    person_id        UUID NOT NULL, -- external_ref: core.persons.id
    graduation_date  DATE,
    alumni_number    VARCHAR(100) UNIQUE,
    status           VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_student_alumni UNIQUE (student_id)
);

CREATE INDEX idx_alumni_student ON alumni (student_id);
CREATE INDEX idx_alumni_person ON alumni (person_id);
