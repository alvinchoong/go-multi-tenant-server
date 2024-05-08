package router

import (
	"context"
	"net/http"
	"strings"
)

// ctxKey is a custom type for context keys
type ctxKey string

// SlugCtxKey is the context key used for storing the slug (subdomain) information
var SlugCtxKey ctxKey = "slug"

// SlugFromCtx extract slug from the context
func SlugFromCtx(ctx context.Context) string {
	if v := ctx.Value(SlugCtxKey); v != nil {
		return v.(string)
	}
	return ""
}

// slugMiddleware extracts the subdomain (tenant identifier) and adds it to the request context
func slugMiddleware(host string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract the subdomain (slug) from the host by removing the primary domain
			subdomain := strings.TrimSuffix(r.Host, "."+host)
			if subdomain != "" && subdomain != host {
				// Store the subdomain in the request context for future use
				ctx = context.WithValue(ctx, SlugCtxKey, subdomain)
			}

			// Serve the request with the modified context
			next.ServeHTTP(w, r.WithContext(ctx))
		})

		return http.HandlerFunc(fn)
	}
}

// slugHandler wraps a custom handler into http.HandlerFunc
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
