-- name: CreateUser :one
INSERT INTO users (id, first_name,last_name,email,password)
VALUES ($1, $2, $3, $4,$5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 OR email = $2;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1 OR email = $2;