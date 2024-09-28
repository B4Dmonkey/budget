-- migrate:up
CREATE TABLE documents_meta (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    persisted_loc VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- migrate:down
DROP TABLE IF EXISTS documents;