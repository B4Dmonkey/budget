CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE documents_meta (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL,
    persisted_loc VARCHAR(255) NOT NULL
  );
CREATE TABLE transactions (
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
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    income_or_expense TEXT CHECK (income_or_expense IN ('income', 'expense')) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL
  );
CREATE TABLE budget_items (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    transaction_date DATE NOT NULL,
    category_id UUID NOT NULL,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (id)
  );
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240928014521'),
  ('20240928155747'),
  ('20240929005504');
