package db

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conns struct {
	Pooled *pgxpool.Pool
	Silos  map[string]*pgxpool.Pool
}

// Connect takes in a pooled connString and a set of connString
func Connect(ctx context.Context, pooledConnString string, silosConnString map[string]string) (*Conns, error) {
	pooled, err := connect(ctx, pooledConnString)
	if err != nil {
		return nil, fmt.Errorf("connect pooled: %w", err)
	}

	silos := make(map[string]*pgxpool.Pool, len(silosConnString))
	for k, v := range silosConnString {
		conn, err := connect(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("connect secondary (%s): %w", k, err)
		}
		silos[k] = conn
	}

	return &Conns{
		Pooled: pooled,
		Silos:  silos,
	}, nil
}

func connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		if v := ctx.Value(SlugKey); v != nil {
			if slug, ok := v.(string); ok && isValidSlug(slug) {
				// set the tenant id into this connection's setting
				_, err := conn.Exec(ctx, "SET app.current_tenant TO '"+slug+"'")
				if err != nil {
					// log the error, and then `return false` to destroy this connection instead of leaving it open.
					slog.Error("BeforeAcquire conn.Exec", slog.Any("err", err))
					return false
				}
			}
		}
		return true
	}
	cfg.AfterRelease = func(conn *pgx.Conn) bool {
		// set the setting to be empty before this connection is released to pool
		_, err := conn.Exec(context.Background(), "RESET app.current_tenant")
		if err != nil {
			// log the error, and then`return false` to destroy this connection instead of leaving it open.
			slog.Error("AfterRelease conn.Exec", slog.Any("err", err))
			return false
		}
		return true
	}

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error("pgxpool.NewWithConfig",
			slog.Any("err", err))
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		slog.Error("conn.Ping",
			slog.Any("err", err))
		return nil, fmt.Errorf("conn.Ping: %w", err)
	}

	slog.Info("connected to database...")
	return conn, nil
}

type key string

var SlugKey key = "slug"

func isValidSlug(slug string) bool {
	// Regular expression to match only alphanumeric characters and underscores
	r := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return r.MatchString(slug)
}
