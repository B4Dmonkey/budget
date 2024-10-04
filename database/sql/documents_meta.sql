-- name: CreateDocumentMeta :one
INSERT INTO
  documents_meta (id, NAME, persisted_loc)
VALUES
  (UUID (), ?, ?) RETURNING *;

-- name: FindOneDocumentMeta :one
SELECT
  id
FROM
  documents_meta
WHERE
  NAME = ?;