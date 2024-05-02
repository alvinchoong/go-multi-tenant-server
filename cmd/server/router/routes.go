package router

import (
	"context"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed static/*
var staticContent embed.FS

func Handler(ctx context.Context, conns *pgxpool.Pool, host string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(extractTenantMiddleware(host))

	// Serve static files from the embedded filesystem
	fs := http.FileServer(http.FS(staticContent))
	r.Get("/static/*", fs.ServeHTTP)

	uh := NewUserHandler(conns)
	r.Post("/api/users", uh.Create())
	r.Get("/api/users", uh.List())
	r.Get("/api/users/{slug}", uh.Get())
	r.Put("/api/users/{slug}", uh.Update())
	r.Delete("/api/users/{slug}", uh.Delete())

	th := NewTodoHandler(conns)
	r.Post("/api/todos", th.Create())
	r.Get("/api/todos", th.List())
	r.Get("/api/todos/{id}", th.Get())
	r.Put("/api/todos/{id}", th.Update())
	r.Patch("/api/todos/{id}", th.Patch())
	r.Delete("/api/todos/{id}", th.Delete())

	return r
}
