-- migrate:up
CREATE TABLE IF NOT EXISTS
  categories (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    income_or_expense TEXT CHECK (income_or_expense IN ('income', 'expense')) NOT NULL,
    amount INTEGER NOT NULL
  );

CREATE TABLE IF NOT EXISTS
  budget_items (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    transaction_date DATE NOT NULL,
    category_id UUID NOT NULL,
    description TEXT NOT NULL,
    amount INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (id)
  );

-- migrate:down
DROP TABLE IF EXISTS budget_items;

DROP TABLE IF EXISTS categories;