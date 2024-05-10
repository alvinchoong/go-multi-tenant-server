-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (slug, description) 
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE slug = $1;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE slug = $1;

-- name: UpdateUser :one
UPDATE users 
SET 
  description = $2,
  updated_at  = now()
WHERE 
  slug = $1
RETURNING *;
