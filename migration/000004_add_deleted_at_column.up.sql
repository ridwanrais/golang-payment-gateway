-- Add deleted_at columns to both tables, which will be NULL if the record is not deleted

ALTER TABLE transactions ADD COLUMN deleted_at TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE virtual_account_transactions ADD COLUMN deleted_at TIMESTAMP WITHOUT TIME ZONE;
