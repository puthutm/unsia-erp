-- pmb_db: 000002_create_applicant_biodata.up.sql

CREATE TABLE applicant_biodata (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id      UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    full_name         VARCHAR(255),
    email             VARCHAR(255),
    phone             VARCHAR(50),
    nik               VARCHAR(50),
    birth_place       VARCHAR(255),
    birth_date        DATE,
    gender            VARCHAR(20),
    religion_id       UUID, -- external_ref: ref.religions.id
    marital_status    VARCHAR(50),
    citizenship       VARCHAR(100),
    jacket_size       VARCHAR(20),
    core_sync_status  VARCHAR(50) DEFAULT 'pending',
    core_synced_at    TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_applicant_biodata_applicant UNIQUE (applicant_id)
);
