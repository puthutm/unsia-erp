-- finance_db: 000016_create_tax_records.up.sql

CREATE TABLE tax_records (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tax_type     VARCHAR(50) NOT NULL,
    source_type  VARCHAR(100),
    source_id    UUID,
    amount       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status       VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    tax_period   DATE NOT NULL
);
