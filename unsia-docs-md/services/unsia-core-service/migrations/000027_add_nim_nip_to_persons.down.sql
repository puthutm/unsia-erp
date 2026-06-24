-- Remove NIM and NIP columns from persons table
-- core_db

DROP INDEX IF EXISTS idx_persons_nim;
DROP INDEX IF EXISTS idx_persons_nip;
DROP INDEX IF EXISTS idx_persons_email;

ALTER TABLE persons DROP COLUMN IF EXISTS nim;
ALTER TABLE persons DROP COLUMN IF EXISTS nip;
