package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"multi-tenant-server/cmd/server/router"
	"multi-tenant-server/internal/db"
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
	var slugDBCfg map[string]string
	if s := os.Getenv("TENANT_DB"); len(s) > 0 {
		if err := json.Unmarshal([]byte(s), &slugDBCfg); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	var silosDB map[string]string
	if s := os.Getenv("DATABASE_SILO_RW_URLS"); len(s) > 0 {
		if err := json.Unmarshal([]byte(s), &silosDB); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	conns, err := db.Connect(ctx, os.Getenv("DATABASE_POOL_RW_URL"), silosDB)
	if err != nil {
		return fmt.Errorf("db.Connect: %w", err)
	}

	// start the server
	slog.Info("starting http server at :8080")
	if err := http.ListenAndServe(":8080", router.Handler(ctx, conns, slugDBCfg)); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}

	return nil
}
