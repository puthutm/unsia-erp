-- pmb_db: 000011_create_loa_documents.up.sql

CREATE TABLE loa_documents (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id  UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    loa_number    VARCHAR(100) UNIQUE NOT NULL,
    file_url      TEXT,
    issued_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    issued_by     UUID, -- external_ref: core.users.id
    CONSTRAINT uq_loa_document_applicant UNIQUE (applicant_id)
);
