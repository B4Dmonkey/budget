CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE documents_meta (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    persisted_loc VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240928014521');
