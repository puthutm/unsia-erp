-- reference_db: 000018_create_status_codes.up.sql

CREATE TABLE status_codes (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module       VARCHAR(50) NOT NULL,
    code         VARCHAR(50) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_status_codes_module_code UNIQUE (module, code)
);

CREATE INDEX idx_status_codes_module ON status_codes(module);
CREATE INDEX idx_status_codes_code ON status_codes(code);
