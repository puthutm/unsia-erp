-- finance_db: 000004_create_payment_gateway_callbacks.up.sql

CREATE TABLE payment_gateway_callbacks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id          UUID REFERENCES payments(id) ON DELETE SET NULL,
    provider            VARCHAR(100) NOT NULL,
    provider_event_id   VARCHAR(255),
    external_reference  VARCHAR(255),
    idempotency_key     VARCHAR(500) UNIQUE,
    payload             JSONB,
    signature_valid     BOOLEAN NOT NULL DEFAULT FALSE,
    callback_status     VARCHAR(50) NOT NULL DEFAULT 'received',
    received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at        TIMESTAMPTZ,
    CONSTRAINT uq_provider_event UNIQUE (provider, provider_event_id)
);

CREATE INDEX idx_callbacks_payment_id ON payment_gateway_callbacks (payment_id);
CREATE INDEX idx_callbacks_external_ref ON payment_gateway_callbacks (external_reference);
