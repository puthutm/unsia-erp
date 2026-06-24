-- pmb_db: 000008_create_applicant_documents.up.sql

CREATE TABLE applicant_documents (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id          UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    document_type_id      UUID NOT NULL, -- external_ref: ref.document_types.id
    file_url              TEXT,
    verification_status   VARCHAR(50) NOT NULL DEFAULT 'pending',
    verification_note     TEXT,
    verified_by           UUID, -- external_ref: core.users.id
    verified_at           TIMESTAMPTZ,
    uploaded_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
