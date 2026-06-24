-- hris_db: 000003_create_functional_positions.up.sql

CREATE TABLE functional_positions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code       VARCHAR(50) UNIQUE NOT NULL,
    name       VARCHAR(255) NOT NULL,
    rank       VARCHAR(50),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE
);
