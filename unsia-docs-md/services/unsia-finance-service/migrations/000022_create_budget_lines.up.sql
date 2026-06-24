-- finance_db: 000022_create_budget_lines.up.sql

CREATE TABLE budget_lines (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id        UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    coa_account_id   UUID REFERENCES coa_accounts(id) ON DELETE SET NULL,
    description      TEXT,
    amount           NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    realized_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00
);
