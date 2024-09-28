-- name: CreateDocumentMeta :one
INSERT INTO documents_meta (name, persisted_loc) VALUES ($1, $2) RETURNING *;