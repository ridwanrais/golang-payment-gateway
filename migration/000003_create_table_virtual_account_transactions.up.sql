-- Specific table for virtual account transactions
CREATE TABLE virtual_account_transactions (
    va_transaction_id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    transaction_id INT REFERENCES transactions(transaction_id),
    bank_name VARCHAR(100) NOT NULL, -- e.g. 'BRI', 'Mandiri', 'BNI'
    virtual_account_number VARCHAR(255) NOT NULL,
    expiry_date TIMESTAMP,
    metadata JSONB, -- for any bank-specific data in JSON format
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for the virtual_account_transactions table
CREATE INDEX idx_va_transactions_transaction_id ON virtual_account_transactions(transaction_id);
CREATE INDEX idx_va_transactions_virtual_account_number ON virtual_account_transactions(virtual_account_number);