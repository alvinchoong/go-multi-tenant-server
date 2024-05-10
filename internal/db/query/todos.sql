-- name: ListTodos :many
SELECT * FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (title, description, user_slug, completed) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1;

-- name: DeleteTodo :execresult
DELETE FROM todos
WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos SET 
  title       = $2,
  description = $3,
  completed   = $4,
  user_slug   = $5,
  updated_at  = now()
WHERE 
  id = $1
RETURNING *;

-- name: PatchTodo :one
UPDATE todos SET 
  title       = COALESCE(sqlc.narg('title'), title),
  description = COALESCE(sqlc.narg('description'), description),
  completed   = COALESCE(sqlc.narg('completed'), completed),
  user_slug   = COALESCE(sqlc.narg('user_slug'), user_slug),
  updated_at  = now()
WHERE
  id = $1
RETURNING *;
