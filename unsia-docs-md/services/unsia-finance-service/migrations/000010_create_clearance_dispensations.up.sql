-- finance_db: 000010_create_clearance_dispensations.up.sql

CREATE TABLE clearance_dispensations (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_clearance_id  UUID NOT NULL REFERENCES student_clearances(id) ON DELETE CASCADE,
    reason                TEXT,
    approved_by           UUID, -- external_ref: core.users.id
    approved_at           TIMESTAMPTZ,
    valid_until           DATE,
    status                VARCHAR(50) NOT NULL DEFAULT 'pending'
);
