-- core_db: 000016_create_integration_event_logs.up.sql

CREATE TABLE integration_event_logs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_module       VARCHAR(50) NOT NULL,
    target_module       VARCHAR(50) NOT NULL,
    event_type          VARCHAR(100) NOT NULL,
    event_key           VARCHAR(500) NOT NULL,
    idempotency_key     VARCHAR(500),
    correlation_id      VARCHAR(100),
    payload             JSONB,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'success', 'failed', 'ignored')),
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at        TIMESTAMPTZ,

    CONSTRAINT uq_integration_event UNIQUE (source_module, target_module, event_type, event_key)
);

CREATE INDEX idx_integration_event_logs_idempotency_key ON integration_event_logs (idempotency_key);
CREATE INDEX idx_integration_event_logs_correlation_id ON integration_event_logs (correlation_id);
CREATE INDEX idx_integration_event_logs_status ON integration_event_logs (status);
CREATE INDEX idx_integration_event_logs_created_at ON integration_event_logs (created_at);
