-- name: CreateTransaction :exec
INSERT INTO
  transactions (
    id,
    details,
    posting_date,
    description,
    amount,
    type,
    balance,
    document_id
  )
VALUES
  (UUID (), ?, ?, ?, ?, ?, ?, ?);