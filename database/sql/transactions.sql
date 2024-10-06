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

-- name: GetPendingTransactions :many
SELECT
  *
FROM
  transactions
WHERE
  document_id = ?
  AND balance IS NULL;

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