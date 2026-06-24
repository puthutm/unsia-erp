-- pmb_db: 000012_create_handover_logs.up.sql

CREATE TABLE handover_logs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_id     UUID NOT NULL REFERENCES applicants(id) ON DELETE CASCADE,
    target_module    VARCHAR(50) NOT NULL DEFAULT 'academic',
    handover_status  VARCHAR(50) NOT NULL DEFAULT 'pending',
    idempotency_key  VARCHAR(500) NOT NULL,
    correlation_id   VARCHAR(100),
    payload          JSONB,
    response_json    JSONB,
    error_message    TEXT,
    handed_over_by   UUID, -- external_ref: core.users.id
    handed_over_at   TIMESTAMPTZ,
    CONSTRAINT uq_handover_logs_idempotency_key UNIQUE (idempotency_key)
);

CREATE INDEX idx_handover_logs_applicant_id ON handover_logs (applicant_id);
CREATE INDEX idx_handover_logs_correlation_id ON handover_logs (correlation_id);
CREATE INDEX idx_handover_logs_status ON handover_logs (handover_status);
