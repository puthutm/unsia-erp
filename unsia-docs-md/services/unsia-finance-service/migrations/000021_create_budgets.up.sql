-- finance_db: 000021_create_budgets.up.sql

CREATE TABLE budgets (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_code   VARCHAR(50) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    fiscal_year   VARCHAR(10) NOT NULL,
    total_amount  NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status        VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
