-- finance_db: 000019_create_journals.up.sql

CREATE TABLE journals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_number  VARCHAR(100) UNIQUE NOT NULL,
    journal_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    source_type     VARCHAR(100),
    source_id       UUID,
    description     TEXT,
    created_by      UUID, -- external_ref: core.users.id
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
