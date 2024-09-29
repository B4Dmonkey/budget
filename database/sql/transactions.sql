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

-- name: GetTransactionsInDateRange :many
SELECT
  *
FROM
  transactions
WHERE
  posting_date >= ?
  AND posting_date <= ?
ORDER BY
  posting_date DESC;