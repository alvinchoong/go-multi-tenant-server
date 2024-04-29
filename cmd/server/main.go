package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"multi-tenant-server/cmd/server/router"
	"multi-tenant-server/internal/db"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
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
	// slog config
	var level slog.Level // default to INFO
	_ = level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))

	// connect to db
	var silosDB map[string]string
	if s := os.Getenv("DATABASE_SILO_RW_URLS"); len(s) > 0 {
		if err := json.Unmarshal([]byte(s), &silosDB); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	beforeAcquire := func(ctx context.Context, conn *pgx.Conn) bool {
		if s := router.SlugFromCtx(ctx); s != "" {
			// set the user for the current session
			rows, err := conn.Query(ctx, "SELECT set_config('app.current_user', $1, false)", s)
			if err != nil {
				// log the error, and then `return false` to destroy this connection instead of leaving it open.
				slog.Error("beforeAcquire conn.Query", slog.Any("err", err))
				return false
			}
			rows.Close()
		}
		return true
	}

	conns, err := db.Connect(ctx, os.Getenv("DATABASE_POOL_RW_URL"), silosDB, beforeAcquire)
	if err != nil {
		return fmt.Errorf("db.Connect: %w", err)
	}

	host := os.Getenv("APP_HOST")
	server := &http.Server{
		Addr:    host,
		Handler: router.Handler(ctx, conns, host),
	}

	// start the server
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gctx.Done()
		slog.Info("server shutting down..")

		defer server.Close()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		slog.Info("server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	return g.Wait()
}
