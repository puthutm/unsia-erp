-- finance_db: 000008_create_clearance_policies.up.sql

CREATE TABLE clearance_policies (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code           VARCHAR(50) UNIQUE NOT NULL,
    name           VARCHAR(255) NOT NULL,
    service_scope  VARCHAR(100),
    rule_json      JSONB,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
