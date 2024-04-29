package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"multi-tenant-server/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TodoHandler struct {
	queries *db.Queries
	conns   *db.Conns
}

func NewTodoHandler(conns *db.Conns) TodoHandler {
	return TodoHandler{
		queries: db.New(),
		conns:   conns,
	}
}

func (h TodoHandler) Create() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		var todo db.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		todo, err = h.queries.CreateTodo(ctx, conn, db.CreateTodoParams{
			Title:       todo.Title,
			Description: todo.Description,
			UserSlug:    SlugFromCtx(ctx),
		})
		if err != nil {
			return fmt.Errorf("create todo failed: %w", err)
		}

		b, err := json.Marshal(todo)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h TodoHandler) List() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		todos, err := h.queries.ListTodos(ctx, conn)
		if err != nil {
			return fmt.Errorf("List todos failed: %w", err)
		}

		b, err := json.Marshal(todos)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}

func (h TodoHandler) Delete() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		res, err := h.queries.DeleteTodo(ctx, conn, id)
		if err != nil {
			return fmt.Errorf("delete todo failed: %w", err)
		}
		if res.RowsAffected() == 0 {
			return fmt.Errorf("todo not found")
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	})
}

func (h TodoHandler) Get() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()
		conn := h.conns.Get(slug)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		todo, err := h.queries.GetTodo(ctx, conn, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("todo not found")
			}
			return fmt.Errorf("get todo failed: %w", err)
		}

		b, err := json.Marshal(todo)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %w", err)
		}

		w.Write(b)

		return nil
	})
}
