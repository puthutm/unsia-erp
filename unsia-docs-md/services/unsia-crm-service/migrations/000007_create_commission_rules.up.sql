-- crm_db: 000007_create_commission_rules.up.sql

CREATE TABLE commission_rules (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referral_type     VARCHAR(50) NOT NULL,
    amount            NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    calculation_type  VARCHAR(50) NOT NULL DEFAULT 'fixed',
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
