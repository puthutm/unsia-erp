-- finance_db: 000020_create_journal_entries.up.sql

CREATE TABLE journal_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_id      UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
    coa_account_id  UUID REFERENCES coa_accounts(id) ON DELETE SET NULL,
    debit           NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    credit          NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    description     TEXT
);
