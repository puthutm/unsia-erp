-- core_db: 000018_create_event_consumers.up.sql

CREATE TABLE event_consumers (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_contract_id       UUID NOT NULL REFERENCES event_contracts(id),
    consumer_module         VARCHAR(50) NOT NULL,
    handler_name            VARCHAR(200),
    retry_policy            JSONB,
    dlq_enabled             BOOLEAN NOT NULL DEFAULT TRUE,
    max_retry               INT NOT NULL DEFAULT 10,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_event_consumer UNIQUE (event_contract_id, consumer_module)
);

CREATE INDEX idx_event_consumers_consumer_module ON event_consumers (consumer_module);
CREATE INDEX idx_event_consumers_is_active ON event_consumers (is_active);
