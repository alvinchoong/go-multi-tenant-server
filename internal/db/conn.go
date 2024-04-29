package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conns struct {
	Pooled *pgxpool.Pool
	Silos  map[string]*pgxpool.Pool
}

// Get the corresponding db conn for tenant, default to "pooled" db
func (c *Conns) Get(slug string) *pgxpool.Pool {
	if conn, ok := c.Silos[slug]; ok {
		return conn
	}

	return c.Pooled
}

// Connect takes in a pooled connString and a set of connString
func Connect(
	ctx context.Context,
	pooledConnString string,
	silosConnString map[string]string,
	beforeAccquire func(context.Context, *pgx.Conn) bool,
) (*Conns, error) {
	pooled, err := connect(ctx, pooledConnString, beforeAccquire)
	if err != nil {
		return nil, fmt.Errorf("connect pooled: %w", err)
	}
	slog.Info("connected to database (pooled)")

	silos := make(map[string]*pgxpool.Pool, len(silosConnString))
	for k, v := range silosConnString {
		conn, err := connect(ctx, v, beforeAccquire)
		if err != nil {
			return nil, fmt.Errorf("connect secondary (%s): %w", k, err)
		}
		silos[k] = conn
		slog.Info("connected to database (" + k + ")")
	}

	return &Conns{
		Pooled: pooled,
		Silos:  silos,
	}, nil
}

func connect(
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
