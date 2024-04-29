-- name: ListTodos :many
SELECT * FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (title, description, user_slug) 
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1;

-- name: DeleteTodo :execresult
DELETE FROM todos
WHERE id = $1;

-- name: UpdateTodo :execresult
UPDATE todos SET 
  title = $2,
  description = $3,
  completed = $4
WHERE id = $1;
