package router

import (
	"context"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
)

func Handler(ctx context.Context, conns *db.Conns, host string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(slugMiddleware(host))

	uh := NewUserHandler(conns)
	r.Post("/api/users", uh.Create())
	r.Get("/api/users", uh.List())
	r.Delete("/api/users/{slug}", uh.Delete())
	r.Get("/api/users/{slug}", uh.Get())

	th := NewTodoHandler(conns)
	r.Post("/api/todos", th.Create())
	r.Get("/api/todos", th.List())
	r.Delete("/api/todos/{id}", th.Delete())
	r.Get("/api/todos/{id}", th.Get())

	return r
}
