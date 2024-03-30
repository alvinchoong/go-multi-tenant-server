package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	// start the server
	slog.Info("starting http server at :8080")
	if err := http.ListenAndServe(":8080", Handler(ctx)); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}

	return nil
}

func Handler(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Use(TenantSlugMiddleware)

	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		slug, ok := ctx.Value(slugKey).(string)
		if !ok {
			fmt.Fprintf(w, "no slug")
			return
		}

		fmt.Fprintf(w, "slug: %s", slug)
	})

	return r
}

type key string

var slugKey key = "slug"

func TenantSlugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract subdomain from request and set it in the context (not the best way)
		parts := strings.Split(r.Host, ".")
		slug := parts[0]

		ctx := context.WithValue(r.Context(), slugKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
