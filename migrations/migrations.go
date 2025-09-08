package migrations

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func Run(ctx context.Context, connStr string) error {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	path := "migrations/001_init.sql"

	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	sql := string(sqlBytes)

	_, err = pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
