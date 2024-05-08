package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	queries *db.Queries
	conn    db.DBTX
}

func NewUserHandler(conn db.DBTX) UserHandler {
	return UserHandler{
		queries: db.New(),
		conn:    conn,
	}
}

func (h UserHandler) Create() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		user, err = h.queries.CreateUser(ctx, h.conn, db.CreateUserParams{
			Slug:        user.Slug,
			Description: user.Description,
		})
		if err != nil {
			return fmt.Errorf("CreateUser failed: %w", err)
		}

		b, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h UserHandler) List() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		users, err := h.queries.ListUsers(ctx, h.conn)
		if err != nil {
			return fmt.Errorf("ListUsers failed: %w", err)
		}

		b, err := json.Marshal(users)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h UserHandler) Delete() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		res, err := h.queries.DeleteUser(ctx, h.conn, chi.URLParam(r, "slug"))
		if err != nil {
			return fmt.Errorf("DeleteUser failed: %w", err)
		}
		if res.RowsAffected() == 0 {
			return errors.New("user not found")
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	})
}

func (h UserHandler) Get() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		user, err := h.queries.GetUser(ctx, h.conn, chi.URLParam(r, "slug"))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("GetUser failed: %w", err)
		}

		b, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h UserHandler) Update() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		user, err = h.queries.UpdateUser(ctx, h.conn, db.UpdateUserParams{
			Slug:        chi.URLParam(r, "slug"),
			Description: user.Description,
		})
		if err != nil {
			return fmt.Errorf("UpdateUser failed: %w", err)
		}

		b, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}
