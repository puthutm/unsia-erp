-- academic_db: 000004_create_nim_sequences.up.sql

CREATE TABLE nim_sequences (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    study_program_id  UUID NOT NULL, -- external_ref: ref.study_programs.id
    entry_period_id   UUID NOT NULL, -- external_ref: ref.academic_periods.id
    sequence_year     VARCHAR(10) NOT NULL,
    last_number       INT NOT NULL DEFAULT 0,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_nim_sequence UNIQUE (study_program_id, entry_period_id, sequence_year)
);

CREATE INDEX idx_nim_seq_study_program ON nim_sequences (study_program_id);
CREATE INDEX idx_nim_seq_entry_period ON nim_sequences (entry_period_id);
