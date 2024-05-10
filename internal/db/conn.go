package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect takes in a connection string and returns a pgxpool
func Connect(
	ctx context.Context,
	connString string,
	beforeAccquire func(context.Context, *pgx.Conn) bool,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	cfg.BeforeAcquire = beforeAccquire

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

	return conn, nil
}
