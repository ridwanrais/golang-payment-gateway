CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- General transaction table
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    reference_number VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    transaction_date TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    transaction_amount NUMERIC(15,2) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL, -- e.g. 'completed', 'pending', 'failed'
    transaction_type VARCHAR(50) NOT NULL, -- e.g. 'virtual_account', 'credit_card', 'bank_transfer'
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for the transactions table
CREATE INDEX idx_transactions_reference_number ON transactions(reference_number);

