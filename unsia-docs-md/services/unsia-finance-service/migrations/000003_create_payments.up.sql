-- finance_db: 000003_create_payments.up.sql

CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id          UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    payment_method_id   UUID, -- external_ref: ref.payment_methods.id
    payment_number      VARCHAR(100) UNIQUE,
    amount              NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    payment_status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    paid_at             TIMESTAMPTZ,
    external_reference  VARCHAR(255),
    idempotency_key     VARCHAR(500) UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_invoice_id ON payments (invoice_id);
CREATE INDEX idx_payments_payment_status ON payments (payment_status);
CREATE INDEX idx_payments_external_reference ON payments (external_reference);
CREATE INDEX idx_payments_invoice_status ON payments (invoice_id, payment_status);
