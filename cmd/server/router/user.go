package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID         string    `json:"id"`
	TenantSlug string    `json:"tenant_slug"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func userCreate(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "json.Decode failed", http.StatusBadRequest)
			return
		}

		err = db.QueryRow(ctx, `
		INSERT INTO users (id, tenant_slug) 
		VALUES ($1, $2) RETURNING *`, u.ID, u.TenantSlug).
			Scan(&u.ID, &u.TenantSlug, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.Exec failed: %+v", err), http.StatusInternalServerError)
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

func userList(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		rows, err := db.Query(ctx, "SELECT * FROM users")
		if err != nil {
			http.Error(w, fmt.Sprintf("db.query failed: %+v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
		if err != nil {
			http.Error(w, fmt.Sprintf("pgx.CollectRows failed: %+v", err), http.StatusInternalServerError)
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

func userDelete(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		id := chi.URLParam(r, "id")

		_, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.query failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func userGet(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		id := chi.URLParam(r, "id")

		_, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.query failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
