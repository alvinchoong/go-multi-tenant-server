package router

import (
	"context"
	"net/http"
	"strings"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Handler(ctx context.Context, conns *db.Conns, slugDBCfg map[string]string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(tenantSlugMiddleware)

	r.Post("/api/tenants", tenantCreate(conns, slugDBCfg))
	r.Get("/api/tenants", tenantList(conns, slugDBCfg))

	r.Post("/api/users", userCreate(conns, slugDBCfg))
	r.Get("/api/users", userList(conns, slugDBCfg))
	r.Delete("/api/users/{id}", userDelete(conns, slugDBCfg))
	r.Get("/api/users/{id}", userGet(conns, slugDBCfg))

	return r
}

// tenantSlugMiddleware extract subdomain from request and set it in the context
func tenantSlugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not the best way
		parts := strings.Split(r.Host, ".")
		slug := parts[0]

		ctx := context.WithValue(r.Context(), db.SlugKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// slugFromContext retrieve the slug from context
func slugFromContext(ctx context.Context) string {
	var slug string
	if v := ctx.Value(db.SlugKey); v != nil {
		if s, ok := v.(string); ok {
			slug = s
		}
	}

	return slug
}

// pickDBConn returns the corresponding db conn for tenant, default to "pooled" db
func pickDBConn(p *db.Conns, slugDBCfg map[string]string, slug string) *pgxpool.Pool {
	if cfg, ok := slugDBCfg[slug]; ok {
		if conn, ok := p.Silos[cfg]; ok {
			return conn
		}
	}

	return p.Pooled
}
