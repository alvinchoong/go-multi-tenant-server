package router

import (
	"context"
	"net/http"
	"strings"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
)

func Handler(ctx context.Context, conns *db.Conns, host string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(tenantSlugMiddleware(host))

	r.Post("/api/users", userCreate(conns))
	r.Get("/api/users", userList(conns))
	r.Delete("/api/users/{id}", userDelete(conns))
	r.Get("/api/users/{id}", userGet(conns))

	th := NewTodoHandler(conns)
	r.Post("/api/todos", th.Create())
	r.Get("/api/todos", th.List())
	r.Delete("/api/todos/{id}", th.Delete())
	r.Get("/api/todos/{id}", th.Get())

	return r
}

// tenantSlugMiddleware extract subdomain from request and set it in the context
func tenantSlugMiddleware(host string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			subdomain := strings.TrimSuffix(r.Host, "."+host)
			if subdomain != "" && subdomain != host {
				ctx = context.WithValue(ctx, db.SlugCtxKey, subdomain)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})

		return http.HandlerFunc(fn)
	}
}

func SlugFromCtx(ctx context.Context) string {
	var s string
	if v := ctx.Value(db.SlugCtxKey); v != nil {
		s = v.(string)
	}
	return s
}
