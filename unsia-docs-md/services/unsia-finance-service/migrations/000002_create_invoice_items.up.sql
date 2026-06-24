-- finance_db: 000002_create_invoice_items.up.sql

CREATE TABLE invoice_items (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id            UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    payment_component_id  UUID, -- external_ref: ref.payment_components.id
    description           TEXT,
    amount                NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    discount_amount       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    final_amount          NUMERIC(15,2) NOT NULL DEFAULT 0.00
);
