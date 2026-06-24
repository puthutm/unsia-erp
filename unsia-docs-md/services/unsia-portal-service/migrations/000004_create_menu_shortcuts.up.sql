-- portal_db: 000004_create_menu_shortcuts.up.sql

CREATE TABLE menu_shortcuts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL, -- external_ref: core.users.id
    menu_code   VARCHAR(100) NOT NULL,
    menu_label  VARCHAR(255) NOT NULL,
    target_url  TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    
    CONSTRAINT uq_user_menu_code UNIQUE (user_id, menu_code)
);

CREATE INDEX idx_menu_shortcuts_user ON menu_shortcuts (user_id);
