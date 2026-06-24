-- finance_db: 000007_create_installment_requests.up.sql

CREATE TABLE installment_requests (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id    UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    student_id    UUID, -- external_ref: academic.students.id
    status        VARCHAR(50) NOT NULL DEFAULT 'pending',
    reason        TEXT,
    requested_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_by   UUID, -- external_ref: core.users.id
    approved_at   TIMESTAMPTZ
);
