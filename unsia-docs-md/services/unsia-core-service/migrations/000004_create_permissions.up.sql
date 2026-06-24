-- core_db: 000004_create_permissions.up.sql
-- Pola permission: module.resource.action (contoh: pmb.applicants.create)

CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(200) UNIQUE NOT NULL,
    module      VARCHAR(50) NOT NULL,
    resource    VARCHAR(100) NOT NULL,
    action      VARCHAR(50) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_permissions_module ON permissions (module);
CREATE INDEX idx_permissions_code ON permissions (code);
