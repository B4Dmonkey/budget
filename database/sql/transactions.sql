-- name: CreateTransaction :exec
INSERT INTO transactions (
    id,
    details,
    posting_date,
    description,
    amount,
    type,
    balance,
    document_id
  )
VALUES (UUID (), ?, ?, ?, ?, ?, ?, ?);
-- name: GetIncomeVsExpenseAndTotal :one
SELECT SUM(
    CASE
      WHEN amount > 0 THEN amount
      ELSE 0
    END
  ) AS income,
  SUM(
    CASE
      WHEN amount < 0 THEN amount
      ELSE 0
    END
  ) AS expense,
  SUM(amount) AS TOTAL
FROM transactions;
-- name: GetPendingTransactions :many
SELECT *
FROM transactions
WHERE document_id = ?
  AND balance IS NULL;
-- name: GetTransactionsInDateRange :many
SELECT sqlc.embed(transactions),
  sqlc.embed(documents_meta)
FROM transactions
  JOIN documents_meta ON transactions.document_id = documents_meta.id
WHERE transactions.posting_date >= @start_date
  AND transactions.posting_date <= @end_date
  AND documents_meta.publishing_date =(
    SELECT MAX(publishing_date)
    FROM documents_meta
  )
ORDER BY transactions.posting_date DESC;