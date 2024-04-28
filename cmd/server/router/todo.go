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
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		var todo db.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, "json.Decode failed", http.StatusBadRequest)
			return
		}

		todo, err = h.queries.CreateTodo(ctx, conn, db.CreateTodoParams{
			Title:       todo.Title,
			Description: todo.Description,
			UserSlug:    SlugFromCtx(ctx),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("CreateTodo failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(todo)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (h TodoHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		todos, err := h.queries.ListTodos(ctx, conn)
		if err != nil {
			http.Error(w, fmt.Sprintf("ListTodos failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (h TodoHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid ID: %+v", err), http.StatusBadRequest)
			return
		}

		res, err := h.queries.DeleteTodo(ctx, conn, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("DeleteTodo failed: %+v", err), http.StatusInternalServerError)
			return
		}
		if res.RowsAffected() == 0 {
			http.Error(w, "todo not found", http.StatusNotFound)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h TodoHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := h.conns.Get(ctx)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid ID: %+v", err), http.StatusBadRequest)
			return
		}

		todo, err := h.queries.GetTodo(ctx, conn, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("GetTodo failed: %+v", err), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(todo)
		if err != nil {
			http.Error(w, fmt.Sprintf("json.Marshal failed: %+v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
