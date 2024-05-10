package router

import (
	"context"
	"net/http"
	"strings"
)

// ctxKey is a custom type for context keys
type ctxKey string

// TenantCtxKey is the context key used for storing the tenant identifier (subdomain)
var TenantCtxKey ctxKey = "tenant"

// TenantFromCtx extracts the tenant identifier from the context
func TenantFromCtx(ctx context.Context) string {
	if v := ctx.Value(TenantCtxKey); v != nil {
		return v.(string)
	}
	return ""
}

// extractTenantMiddleware extracts the subdomain (tenant identifier) and adds it to the request context
func extractTenantMiddleware(host string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract the subdomain (tenant identifier) by removing the main domain
			subdomain := strings.TrimSuffix(r.Host, "."+host)
			if subdomain != "" && subdomain != host {
				// Store the subdomain in the request context for future access
				ctx = context.WithValue(ctx, TenantCtxKey, subdomain)
			}

			// Serve the request with the modified context
			next.ServeHTTP(w, r.WithContext(ctx))
		})

		return http.HandlerFunc(fn)
	}
}

// tenantHandler wraps a custom handler into http.HandlerFunc with tenant identifier support
func tenantHandler(fn func(w http.ResponseWriter, r *http.Request, slug string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := TenantFromCtx(ctx)

		w.Header().Set("Content-Type", "application/json")
		if err := fn(w, r, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// tenantUIHandler wraps a custom handler into http.HandlerFunc with tenant identifier support
func tenantUIHandler(fn func(w http.ResponseWriter, r *http.Request, slug string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := TenantFromCtx(ctx)

		w.Header().Set("Content-Type", "text/html")
		if err := fn(w, r, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
