package router

import (
	"context"
	"net/http"
	"strings"
)

// ctxKey is a custom type for context keys
type ctxKey string

// SlugCtxKey is a context key for the slug
var SlugCtxKey ctxKey = "slug"

// SlugFromCtx extract slug from the context
func SlugFromCtx(ctx context.Context) string {
	if v := ctx.Value(SlugCtxKey); v != nil {
		return v.(string)
	}
	return ""
}

// slugMiddleware extract subdomain from request and set it in the context
func slugMiddleware(host string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			subdomain := strings.TrimSuffix(r.Host, "."+host)
			if subdomain != "" && subdomain != host {
				ctx = context.WithValue(ctx, SlugCtxKey, subdomain)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})

		return http.HandlerFunc(fn)
	}
}

func slugHandler(fn func(w http.ResponseWriter, r *http.Request, slug string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := SlugFromCtx(ctx)

		w.Header().Set("Content-Type", "application/json")
		if err := fn(w, r, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
