// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

const CreateUser = `-- name: CreateUser :one
INSERT INTO users (slug) 
VALUES ($1)
RETURNING slug, created_at, updated_at
`

func (q *Queries) CreateUser(ctx context.Context, db DBTX, slug string) (User, error) {
	row := db.QueryRow(ctx, CreateUser, slug)
	var i User
	err := row.Scan(&i.Slug, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const DeleteUser = `-- name: DeleteUser :execresult
DELETE FROM users
WHERE slug = $1
`

func (q *Queries) DeleteUser(ctx context.Context, db DBTX, slug string) (pgconn.CommandTag, error) {
	return db.Exec(ctx, DeleteUser, slug)
}

const GetUser = `-- name: GetUser :one
SELECT slug, created_at, updated_at FROM users 
WHERE slug = $1
`

func (q *Queries) GetUser(ctx context.Context, db DBTX, slug string) (User, error) {
	row := db.QueryRow(ctx, GetUser, slug)
	var i User
	err := row.Scan(&i.Slug, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const ListUsers = `-- name: ListUsers :many
SELECT slug, created_at, updated_at FROM users
`

func (q *Queries) ListUsers(ctx context.Context, db DBTX) ([]User, error) {
	rows, err := db.Query(ctx, ListUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.Slug, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
