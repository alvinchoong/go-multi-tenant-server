package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var slug = "abcdef"

func BenchmarkSetConfigExec(b *testing.B) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close(ctx)
	b.Cleanup(func() {
		conn.Close(ctx)
	})

	for i := 0; i < b.N; i++ {
		_, err := conn.Exec(ctx, "select set_config('app.current_tenant', $1, false)", slug)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetConfigQueryRow(b *testing.B) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close(ctx)
	b.Cleanup(func() {
		conn.Close(ctx)
	})

	for i := 0; i < b.N; i++ {
		var s string
		err := conn.QueryRow(ctx, "select set_config('app.current_tenant', $1, false)", slug).Scan(&s)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetConfigQuery(b *testing.B) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close(ctx)
	b.Cleanup(func() {
		conn.Close(ctx)
	})

	for i := 0; i < b.N; i++ {
		_, err := conn.Query(ctx, "select set_config('app.current_tenant', $1, false)", slug)
		if err != nil {
			b.Fatal(err)
		}
	}
}
