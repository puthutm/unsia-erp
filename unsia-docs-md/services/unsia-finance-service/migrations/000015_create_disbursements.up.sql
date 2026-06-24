-- finance_db: 000015_create_disbursements.up.sql

CREATE TABLE disbursements (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disbursement_type     VARCHAR(50) NOT NULL,
    commission_record_id  UUID, -- external_ref: crm.commission_records.id
    recipient_person_id   UUID, -- external_ref: core.persons.id
    amount                NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status                VARCHAR(50) NOT NULL DEFAULT 'pending',
    disbursed_at          TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
