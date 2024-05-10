package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"multi-tenant-server/cmd/server/router/components"
	"multi-tenant-server/internal/db"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type TodoUIHandler struct {
	queries *db.Queries
	conn    db.DBTX
}

func NewTodoUIHandler(conn db.DBTX) TodoUIHandler {
	return TodoUIHandler{
		queries: db.New(),
		conn:    conn,
	}
}

func (h TodoUIHandler) Index() http.HandlerFunc {
	return tenantUIHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		todos, err := h.queries.ListTodos(ctx, h.conn)
		if err != nil {
			return fmt.Errorf("ListTodos failed: %w", err)
		}

		if err := components.Layout(components.TodoIndex(todos)).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) New() http.HandlerFunc {
	return tenantUIHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		if err := components.TodoModal(components.TodoModalParams{
			HXMethod: components.HXMethodPost,
			Action:   "/todos",
			Heading:  "Create Todo",
			BtnText:  "Submit",
			FormAttrs: templ.Attributes{
				"hx-target": "#table-body",
				"hx-swap":   "beforeend",
			},
		}).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) Create() http.HandlerFunc {
	return tenantUIHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()

		p := db.CreateTodoParams{
			UserSlug:    slug,
			Title:       postFormValue(r, "title"),
			Description: nil,
			Completed:   postFormValue(r, "completed") == "on",
		}

		if s := postFormValue(r, "description"); s != "" {
			p.Description = &s
		}

		todo, err := h.queries.CreateTodo(ctx, h.conn, p)
		if err != nil {
			return fmt.Errorf("create todo failed: %w", err)
		}

		w.Header().Set("hx-trigger-after-settle", "closeModal")
		if err := components.TodoRow(todo).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) Get() http.HandlerFunc {
	return tenantUIHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
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

		if err := components.TodoModal(components.TodoModalParams{
			HXMethod: components.HXMethodPut,
			Action:   "/todos/" + idStr,
			Heading:  "Update Todo",
			BtnText:  "Submit",
			FormAttrs: templ.Attributes{
				"hx-target": "#table-row-" + idStr,
				"hx-swap":   "outerHTML",
			},
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) Update() http.HandlerFunc {
	return tenantHandler(func(w http.ResponseWriter, r *http.Request, slug string) error {
		ctx := r.Context()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		p := db.UpdateTodoParams{
			ID:          id,
			UserSlug:    slug,
			Title:       postFormValue(r, "title"),
			Description: nil,
			Completed:   postFormValue(r, "completed") == "on",
		}
		if s := postFormValue(r, "description"); s != "" {
			p.Description = &s
		}

		todo, err := h.queries.UpdateTodo(ctx, h.conn, p)
		if err != nil {
			return fmt.Errorf("create todo failed: %w", err)
		}

		w.Header().Set("hx-trigger-after-settle", "closeModal")
		if err := components.TodoRow(todo).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) Patch() http.HandlerFunc {
	return tenantHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return fmt.Errorf("invalid todo id: %w", err)
		}

		completed := postFormValue(r, "completed") == "on"

		p := db.PatchTodoParams{
			ID:        id,
			Completed: &completed,
		}

		todo, err := h.queries.PatchTodo(ctx, h.conn, p)
		if err != nil {
			return fmt.Errorf("create todo failed: %w", err)
		}

		if err := components.TodoRow(todo).Render(ctx, w); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		return nil
	})
}

func (h TodoUIHandler) Destroy() http.HandlerFunc {
	return tenantUIHandler(func(w http.ResponseWriter, r *http.Request, _ string) error {
		ctx := r.Context()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

		return nil
	})
}

func postFormValue(r *http.Request, key string) string {
	return strings.TrimSpace(r.PostFormValue(key))
}
