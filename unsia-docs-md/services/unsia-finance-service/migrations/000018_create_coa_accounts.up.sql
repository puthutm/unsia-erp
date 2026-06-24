-- finance_db: 000018_create_coa_accounts.up.sql

CREATE TABLE coa_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code    VARCHAR(50) UNIQUE NOT NULL,
    account_name    VARCHAR(255) NOT NULL,
    normal_balance  VARCHAR(20) NOT NULL CHECK (normal_balance IN ('DEBIT', 'CREDIT')),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
