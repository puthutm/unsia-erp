-- portal_db: 000007_create_inbox_events.up.sql

CREATE TABLE inbox_events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name          VARCHAR(200) NOT NULL,
    event_version       VARCHAR(10) NOT NULL,
    event_key           VARCHAR(500) NOT NULL,
    publisher_module    VARCHAR(50) NOT NULL,
    publisher_database  VARCHAR(50),
    consumer_module     VARCHAR(50) NOT NULL DEFAULT 'portal',
    aggregate_type      VARCHAR(100),
    aggregate_id        UUID,
    payload             JSONB NOT NULL,
    payload_hash        TEXT,
    headers             JSONB,
    correlation_id      VARCHAR(100),
    causation_id        VARCHAR(100),
    status              VARCHAR(20) NOT NULL DEFAULT 'RECEIVED'
                        CHECK (status IN ('RECEIVED', 'PROCESSED', 'RETRYING', 'FAILED', 'DLQ', 'IGNORED_DUPLICATE')),
    retry_count         INT NOT NULL DEFAULT 0,
    max_retry           INT NOT NULL DEFAULT 10,
    next_retry_at       TIMESTAMPTZ,
    locked_at           TIMESTAMPTZ,
    locked_by           VARCHAR(100),
    last_error          TEXT,
    received_at         TIMESTAMPTZ,
    processed_at        TIMESTAMPTZ,
    dead_letter_at      TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_consumer_event_key UNIQUE (consumer_module, event_key)
);

CREATE INDEX idx_inbox_events_event_name ON inbox_events (event_name);
CREATE INDEX idx_inbox_events_publisher_module ON inbox_events (publisher_module);
CREATE INDEX idx_inbox_events_status ON inbox_events (status);
CREATE INDEX idx_inbox_events_next_retry_at ON inbox_events (next_retry_at);
CREATE INDEX idx_inbox_events_aggregate ON inbox_events (aggregate_type, aggregate_id);
CREATE INDEX idx_inbox_events_correlation_id ON inbox_events (correlation_id);
CREATE INDEX idx_inbox_events_dead_letter_at ON inbox_events (dead_letter_at);
