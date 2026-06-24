-- finance_db: 000026_create_reconciliation_mismatch_logs.up.sql

CREATE TABLE reconciliation_mismatch_logs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_module       VARCHAR(50) NOT NULL,
    source_table        VARCHAR(100) NOT NULL,
    source_ref_id       UUID,
    consumer_module     VARCHAR(50) NOT NULL DEFAULT 'finance',
    consumer_table      VARCHAR(100),
    consumer_ref_id     UUID,
    source_event_key    VARCHAR(500),
    mismatch_type       VARCHAR(30) NOT NULL
                        CHECK (mismatch_type IN ('missing_source', 'missing_snapshot', 'value_mismatch', 'stale_snapshot', 'duplicate_projection')),
    source_value        JSONB,
    snapshot_value      JSONB,
    status              VARCHAR(20) NOT NULL DEFAULT 'OPEN'
                        CHECK (status IN ('OPEN', 'CORRECTED', 'IGNORED', 'PENDING_REVIEW')),
    reason              TEXT,
    detected_at         TIMESTAMPTZ,
    corrected_at        TIMESTAMPTZ,
    ignored_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reconciliation_status ON reconciliation_mismatch_logs (status);
CREATE INDEX idx_reconciliation_mismatch_type ON reconciliation_mismatch_logs (mismatch_type);
CREATE INDEX idx_reconciliation_source_module ON reconciliation_mismatch_logs (source_module);
CREATE INDEX idx_reconciliation_consumer_module ON reconciliation_mismatch_logs (consumer_module);
CREATE INDEX idx_reconciliation_source_event_key ON reconciliation_mismatch_logs (source_event_key);
CREATE INDEX idx_reconciliation_detected_at ON reconciliation_mismatch_logs (detected_at);
