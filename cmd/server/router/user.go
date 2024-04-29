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
	conns   *db.Conns
}

func NewUserHandler(conns *db.Conns) UserHandler {
	return UserHandler{
		queries: db.New(),
		conns:   conns,
	}
}

func (h UserHandler) Create() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		var u db.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		u, err = h.queries.CreateUser(ctx, conn, u.Slug)
		if err != nil {
			return fmt.Errorf("CreateUser failed: %w", err)
		}

		b, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h UserHandler) List() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		users, err := h.queries.ListUsers(ctx, conn)
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
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		res, err := h.queries.DeleteUser(ctx, conn, chi.URLParam(r, "slug"))
		if err != nil {
			return fmt.Errorf("DeleteUser failed: %w", err)
		}
		if res.RowsAffected() == 0 {
			http.Error(w, "user not found", http.StatusNotFound)
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	})
}

func (h UserHandler) Get() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		u, err := h.queries.GetUser(ctx, conn, chi.URLParam(r, "slug"))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("GetUser failed: %w", err)
		}

		b, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}
