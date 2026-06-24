-- hris_db: 000012_create_outbox_events.up.sql

CREATE TABLE outbox_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name      VARCHAR(200) NOT NULL,
    event_version   VARCHAR(10) NOT NULL DEFAULT 'v1',
    event_key       VARCHAR(500) UNIQUE NOT NULL,
    event_type      VARCHAR(30) NOT NULL
                    CHECK (event_type IN ('DOMAIN_EVENT', 'INTEGRATION_EVENT', 'NOTIFICATION_EVENT', 'SNAPSHOT_EVENT')),
    aggregate_type  VARCHAR(100) NOT NULL,
    aggregate_id    UUID NOT NULL,
    payload         JSONB NOT NULL,
    payload_hash    TEXT,
    headers         JSONB,
    idempotency_key VARCHAR(500),
    correlation_id  VARCHAR(100),
    causation_id    VARCHAR(100),
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING'
                    CHECK (status IN ('PENDING', 'PUBLISHED', 'RETRYING', 'FAILED', 'DLQ')),
    retry_count     INT NOT NULL DEFAULT 0,
    max_retry       INT NOT NULL DEFAULT 10,
    next_retry_at   TIMESTAMPTZ,
    locked_at       TIMESTAMPTZ,
    locked_by       VARCHAR(100),
    last_error      TEXT,
    occurred_at     TIMESTAMPTZ NOT NULL,
    published_at    TIMESTAMPTZ,
    dead_letter_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_outbox_events_event_name ON outbox_events (event_name);
CREATE INDEX idx_outbox_events_status ON outbox_events (status);
CREATE INDEX idx_outbox_events_next_retry_at ON outbox_events (next_retry_at);
CREATE INDEX idx_outbox_events_aggregate ON outbox_events (aggregate_type, aggregate_id);
CREATE INDEX idx_outbox_events_correlation_id ON outbox_events (correlation_id);
CREATE INDEX idx_outbox_events_dead_letter_at ON outbox_events (dead_letter_at);
