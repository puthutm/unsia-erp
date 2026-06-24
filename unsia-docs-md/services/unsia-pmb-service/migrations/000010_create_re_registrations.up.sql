-- pmb_db: 000010_create_re_registrations.up.sql

CREATE TABLE re_registrations (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id  UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    status        VARCHAR(50) NOT NULL DEFAULT 'pending',
    submitted_at  TIMESTAMPTZ,
    verified_at   TIMESTAMPTZ,
    verified_by   UUID, -- external_ref: core.users.id
    CONSTRAINT uq_re_registration_applicant UNIQUE (applicant_id)
);
