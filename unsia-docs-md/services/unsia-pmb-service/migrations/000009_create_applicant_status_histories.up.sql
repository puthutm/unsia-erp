-- pmb_db: 000009_create_applicant_status_histories.up.sql

CREATE TABLE applicant_status_histories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id  UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    old_status    VARCHAR(50),
    new_status    VARCHAR(50) NOT NULL,
    changed_by    UUID, -- external_ref: core.users.id
    note          TEXT,
    changed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
