-- academic_db: 000017_create_krs_items.up.sql

CREATE TABLE krs_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    krs_id       UUID NOT NULL REFERENCES krs(id) ON DELETE CASCADE,
    class_id     UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    status       VARCHAR(50) NOT NULL DEFAULT 'selected', -- selected, approved, dropped
    selected_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_krs_class UNIQUE (krs_id, class_id)
);

CREATE INDEX idx_krs_items_krs ON krs_items (krs_id);
CREATE INDEX idx_krs_items_class ON krs_items (class_id);
CREATE INDEX idx_krs_items_status ON krs_items (status);
