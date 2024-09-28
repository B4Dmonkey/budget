-- name: CreateDocumentMeta :one
INSERT INTO documents_meta (id, name, persisted_loc) VALUES (uuid(), $1, $2) RETURNING *;