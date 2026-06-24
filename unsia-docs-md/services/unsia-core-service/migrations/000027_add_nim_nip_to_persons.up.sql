-- Add NIM and NIP columns to persons table
-- core_db

ALTER TABLE persons ADD COLUMN IF NOT EXISTS nim TEXT UNIQUE;
ALTER TABLE persons ADD COLUMN IF NOT EXISTS nip TEXT UNIQUE;

CREATE INDEX IF NOT EXISTS idx_persons_nim ON persons (nim);
CREATE INDEX IF NOT EXISTS idx_persons_nip ON persons (nip);
CREATE INDEX IF NOT EXISTS idx_persons_email ON persons (email);
