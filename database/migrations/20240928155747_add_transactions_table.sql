-- migrate:up
CREATE TABLE
  transactions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    document_id UUID, -- Foreign key column
    details TEXT NOT NULL,
    posting_date DATE NOT NULL,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    type VARCHAR(50) NOT NULL, -- e.g., "credit" or "debit"
    balance DECIMAL(10, 2),
    FOREIGN KEY (document_id) REFERENCES documents_meta (id)
  );

-- migrate:down
DROP TABLE IF EXISTS transactions;