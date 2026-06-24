-- finance_db: 000011_create_cash_accounts.up.sql

CREATE TABLE cash_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code    VARCHAR(50) UNIQUE NOT NULL,
    account_name    VARCHAR(255) NOT NULL,
    bank_name       VARCHAR(100),
    account_number  VARCHAR(100),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
