-- core_db: 000015_create_idempotency_keys.up.sql

CREATE TABLE idempotency_keys (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module              VARCHAR(50) NOT NULL,
    idempotency_key     VARCHAR(500) NOT NULL,
    source_module       VARCHAR(50),
    target_module       VARCHAR(50),
    request_hash        TEXT,
    response_json       JSONB,
    response_payload    JSONB,
    status              VARCHAR(20) NOT NULL DEFAULT 'processing'
                        CHECK (status IN ('processing', 'completed', 'failed', 'expired')),
    locked_until        TIMESTAMPTZ,
    trace_id            VARCHAR(100),
    correlation_id      VARCHAR(100),
    last_error          TEXT,
    expires_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at        TIMESTAMPTZ,

    CONSTRAINT uq_module_idempotency_key UNIQUE (module, idempotency_key)
);

CREATE INDEX idx_idempotency_keys_status ON idempotency_keys (status);
CREATE INDEX idx_idempotency_keys_locked_until ON idempotency_keys (locked_until);
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys (expires_at);
CREATE INDEX idx_idempotency_keys_correlation_id ON idempotency_keys (correlation_id);
