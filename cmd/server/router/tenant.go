package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multi-tenant-server/internal/db"

	"github.com/jackc/pgx/v5"
)

type Tenant struct {
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func tenantCreate(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		var tenant Tenant
		err := json.NewDecoder(r.Body).Decode(&tenant)
		if err != nil {
			http.Error(w, "json.Decode failed", http.StatusBadRequest)
			return
		}

		err = db.QueryRow(ctx, `
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

func tenantList(pool *db.Pool, slugDBCfg map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		rows, err := db.Query(ctx, `SELECT * FROM tenants`)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.Query failed: %+v", err), http.StatusInternalServerError)
			return
		}

		tenants, err := pgx.CollectRows(rows, pgx.RowToStructByName[Tenant])
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
