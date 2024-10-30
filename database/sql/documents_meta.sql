-- name: CreateDocumentMeta :one
INSERT INTO documents_meta (id, NAME, persisted_loc, publishing_date)
VALUES (UUID (), ?, ?, ?)
RETURNING *;
-- name: FindOneDocumentMeta :one
SELECT id
FROM documents_meta
WHERE NAME = ?
  AND publishing_date = ?;