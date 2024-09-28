-- migrate:up
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    details TEXT NOT NULL,
    posting_date DATE NOT NULL,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- e.g., "credit" or "debit"
    balance DECIMAL(10, 2) NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS transactions;
