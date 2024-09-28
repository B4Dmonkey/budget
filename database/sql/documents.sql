-- name: CreateDocument :one
INSERT INTO documents (name, persisted_loc) VALUES ($1, $2) RETURNING *;