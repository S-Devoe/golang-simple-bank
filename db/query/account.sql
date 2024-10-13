-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2;

-- NOTE FOR ME: balance is $2 and id is $1 in the UDEMY course.
-- i want to see what happens if i change the order of the variables in the query. 

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
