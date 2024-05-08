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
	conn    db.DBTX
}

func NewTodoHandler(conn db.DBTX) TodoHandler {
	return TodoHandler{
		queries: db.New(),
		conn:    conn,
	}
}

func (h TodoHandler) Create() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		var todo db.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		todo, err = h.queries.CreateTodo(ctx, h.conn, db.CreateTodoParams{
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			UserSlug:    todo.UserSlug,
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
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		todos, err := h.queries.ListTodos(ctx, h.conn)
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

func (h TodoHandler) Get() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		todo, err := h.queries.GetTodo(ctx, h.conn, id)
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

func (h TodoHandler) Update() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		var todo db.Todo
		err = json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		todo, err = h.queries.UpdateTodo(ctx, h.conn, db.UpdateTodoParams{
			ID:          id,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			UserSlug:    todo.UserSlug,
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

type todoPatchInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
	UserSlug    *string `json:"user_slug"`
}

func (h TodoHandler) Patch() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		var input todoPatchInput
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			return fmt.Errorf("json.Decode failed: %w", err)
		}

		todo, err := h.queries.PatchTodo(ctx, h.conn, db.PatchTodoParams{
			ID:          id,
			Title:       input.Title,
			Description: input.Description,
			Completed:   input.Completed,
			UserSlug:    input.UserSlug,
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

func (h TodoHandler) Delete() http.HandlerFunc {
	return slugHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		res, err := h.queries.DeleteTodo(ctx, h.conn, id)
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
