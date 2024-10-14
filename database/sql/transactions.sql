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
  t.*,
  d.name AS document_name,
  d.publishing_date AS document_publishing_date
FROM
  transactions t
  JOIN documents_meta d ON t.document_id = d.id
WHERE
  t.posting_date >= ?
  AND t.posting_date <= ?
  AND d.publishing_date = (
    SELECT
      MAX(publishing_date)
    FROM
      documents_meta
  )
ORDER BY
  t.posting_date DESC;