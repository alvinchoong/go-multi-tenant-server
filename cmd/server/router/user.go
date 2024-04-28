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
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		var u db.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Decode failed: %+v", err), http.StatusBadRequest)
			return
		}

		u, err = h.queries.CreateUser(ctx, conn, u.Slug)
		if err != nil {
			http.Error(w, fmt.Sprintf("CreateUser failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(u)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (h UserHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		users, err := h.queries.ListUsers(ctx, conn)
		if err != nil {
			http.Error(w, fmt.Sprintf("ListUsers failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(users)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (h UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		slug := chi.URLParam(r, "slug")

		res, err := h.queries.DeleteUser(ctx, conn, slug)
		if err != nil {
			http.Error(w, fmt.Sprintf("DeleteUser failed: %+v", err), http.StatusInternalServerError)
			return
		}
		if res.RowsAffected() == 0 {
			http.Error(w, "user not found", http.StatusNotFound)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h UserHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		slug := chi.URLParam(r, "slug")

		u, err := h.queries.GetUser(ctx, conn, slug)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("GetUser failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(u)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
