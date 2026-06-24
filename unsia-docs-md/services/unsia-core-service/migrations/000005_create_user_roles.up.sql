-- core_db: 000005_create_user_roles.up.sql

CREATE TABLE user_roles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id             UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    study_program_id    UUID,   -- external_ref: reference_db.study_programs.id (nullable; untuk scope admin prodi/kaprodi)
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_user_role_study_program UNIQUE (user_id, role_id, study_program_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);
