package router

import (
	"context"
	"net/http"
	"strings"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Handler(ctx context.Context, pool *db.Pool, slugDBCfg map[string]string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(TenantSlugMiddleware)

	r.Post("/api/tenants", tenantCreate(pool, slugDBCfg))
	r.Get("/api/tenants", tenantList(pool, slugDBCfg))

	r.Post("/api/users", userCreate(pool, slugDBCfg))
	r.Get("/api/users", userList(pool, slugDBCfg))
	r.Delete("/api/users/{id}", userDelete(pool, slugDBCfg))
	r.Get("/api/users/{id}", userGet(pool, slugDBCfg))

	return r
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
