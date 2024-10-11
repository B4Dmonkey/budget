-- migrate:up
ALTER TABLE documents_meta
ADD COLUMN publishing_date TIMESTAMP;

-- migrate:down
ALTER TABLE documents_meta
DROP COLUMN publishing_date;
