-- Drop the deleted_at columns from both tables

ALTER TABLE transactions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE virtual_account_transactions DROP COLUMN IF EXISTS deleted_at;
