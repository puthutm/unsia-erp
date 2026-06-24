-- core_db: 000007_create_applications.up.sql

CREATE TABLE applications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    launch_url  TEXT NOT NULL,
    sso_protocol VARCHAR(50),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);
