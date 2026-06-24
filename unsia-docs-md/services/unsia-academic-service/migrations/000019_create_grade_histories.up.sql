-- academic_db: 000019_create_grade_histories.up.sql

CREATE TABLE grade_histories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grade_id    UUID NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    old_value   JSONB,
    new_value   JSONB,
    changed_by  UUID, -- external_ref: core.users.id
    reason      TEXT,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_grade_histories_grade ON grade_histories (grade_id);
