package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multi-tenant-server/internal/db"

	"github.com/jackc/pgx/v5"
)

type tenant struct {
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func tenantCreate(conns *db.Conns) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := conns.Get(ctx)

		var tenant tenant
		err := json.NewDecoder(r.Body).Decode(&tenant)
		if err != nil {
			http.Error(w, "json.Decode failed", http.StatusBadRequest)
			return
		}

		err = conn.QueryRow(ctx, `
		INSERT INTO tenants (slug, description) 
		VALUES ($1, $2) RETURNING *`, tenant.Slug, tenant.Description).
			Scan(&tenant.Slug, &tenant.Description, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.QueryRow failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(tenant)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func tenantList(conns *db.Conns) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := conns.Get(ctx)

		rows, err := conn.Query(ctx, `SELECT * FROM tenants`)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.Query failed: %+v", err), http.StatusInternalServerError)
			return
		}

		tenants, err := pgx.CollectRows(rows, pgx.RowToStructByName[tenant])
		if err != nil {
			http.Error(w, fmt.Sprintf("pgx.CollectRows failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(tenants)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
