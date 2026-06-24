-- hris_db: 000009_create_performance_reviews.up.sql

CREATE TABLE performance_reviews (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id    UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    review_period  VARCHAR(50) NOT NULL,
    score          NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    status         VARCHAR(50) NOT NULL DEFAULT 'pending',
    reviewed_by    UUID, -- external_ref: core.users.id
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
