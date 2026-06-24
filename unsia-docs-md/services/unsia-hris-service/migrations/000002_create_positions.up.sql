-- hris_db: 000002_create_positions.up.sql

CREATE TABLE positions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code       VARCHAR(50) UNIQUE NOT NULL,
    name       VARCHAR(255) NOT NULL,
    level      VARCHAR(50),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE
);
