-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (slug) 
VALUES ($1)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE slug = $1;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE slug = $1;
