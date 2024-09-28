-- migrate:up
CREATE TABLE
  documents_meta (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL,
    persisted_loc VARCHAR(255) NOT NULL
  );

-- migrate:down
DROP TABLE IF EXISTS documents;