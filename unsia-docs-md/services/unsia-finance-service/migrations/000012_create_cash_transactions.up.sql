-- finance_db: 000012_create_cash_transactions.up.sql

CREATE TABLE cash_transactions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cash_account_id   UUID NOT NULL REFERENCES cash_accounts(id) ON DELETE CASCADE,
    transaction_type  VARCHAR(50) NOT NULL, -- DEBIT, CREDIT
    source_type       VARCHAR(100),
    source_id         UUID,
    amount            NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    description       TEXT,
    transaction_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
