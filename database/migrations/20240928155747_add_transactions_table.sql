-- migrate:up
CREATE TABLE IF NOT EXISTS
  transactions (
    id UUID PRIMARY KEY DEFAULT (UUID()),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    document_id UUID, -- Foreign key column
    details TEXT NOT NULL,
    posting_date DATE NOT NULL,
    description TEXT NOT NULL,
    amount INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL, -- e.g., "credit" or "debit"
    balance INTEGER,
    FOREIGN KEY (document_id) REFERENCES documents_meta (id),
    UNIQUE (posting_date, description, amount, type, balance) -- Composite unique constraint
  );

-- migrate:down
DROP TABLE IF EXISTS transactions;