-- core_db: 000003_create_roles.up.sql

CREATE TABLE roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    scope_type  VARCHAR(30) CHECK (scope_type IN ('global', 'prodi', 'module', 'self')),
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_roles_code ON roles (code);
CREATE INDEX idx_roles_scope_type ON roles (scope_type);
