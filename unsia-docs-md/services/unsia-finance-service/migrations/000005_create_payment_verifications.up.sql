-- finance_db: 000005_create_payment_verifications.up.sql

CREATE TABLE payment_verifications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id          UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    verified_by         UUID, -- external_ref: core.users.id
    verification_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    rejection_reason    VARCHAR(255),
    note                TEXT,
    verified_at         TIMESTAMPTZ
);
