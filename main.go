package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"multi-tenant-server/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) // the only background ctx
	if err := errmain(ctx); err != nil {
		cancel()
		slog.Error("something happen", slog.Any("err", err))
		os.Exit(1)
	}
}

func errmain(ctx context.Context) error {
	pool, err := db.Connect(ctx,
		os.Getenv("DATABASE_PRIMARY_RW_URL"),
		map[string]string{
			"special": os.Getenv("DATABASE_SECONDARY_RW_URL"),
		})
	if err != nil {
		return fmt.Errorf("db.Connect: %w", err)
	}

	slugDBCfg := map[string]string{
		"special-abc": "special",
		"special-def": "special",
	}

	// start the server
	slog.Info("starting http server at :8080")
	if err := http.ListenAndServe(":8080", Handler(ctx, pool, slugDBCfg)); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}

	return nil
}

func Handler(ctx context.Context, pool *db.Pool, slugDBCfg map[string]string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(TenantSlugMiddleware)

	r.Post("/api/tenants", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Get("/api/tenants", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Post("/api/users", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Delete("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Get("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slug := slugFromContext(ctx)
		db := pickDBPool(pool, slugDBCfg, slug)

		id := chi.URLParam(r, "id")

		var u User
		err := db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).
			Scan(&u.ID, &u.TenantSlug, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("db.query failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(u)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	return r
}

type Tenant struct {
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID         string    `json:"id"`
	TenantSlug string    `json:"tenant_slug"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func TenantSlugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract subdomain from request and set it in the context (not the best way)
		parts := strings.Split(r.Host, ".")
		slug := parts[0]

		ctx := context.WithValue(r.Context(), db.SlugKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func slugFromContext(ctx context.Context) string {
	var slug string
	if v := ctx.Value(db.SlugKey); v != nil {
		if s, ok := v.(string); ok {
			slug = s
		}
	}

	return slug
}

func pickDBPool(p *db.Pool, slugDBCfg map[string]string, slug string) *pgxpool.Pool {
	if cfg, ok := slugDBCfg[slug]; ok {
		if conn, ok := p.Secondary[cfg]; ok {
			return conn
		}
	}

	return p.Primary
}
