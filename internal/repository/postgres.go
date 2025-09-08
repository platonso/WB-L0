package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/order-viewer/internal/domain"
)

type PostgresRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, connStr string) (*PostgresRepo, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &PostgresRepo{DB: db}, nil
}

func (r *PostgresRepo) Close() {
	if r.DB != nil {
		r.DB.Close()
	}
}

func (r *PostgresRepo) Save(ctx context.Context, order *domain.Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresRepo) FindByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}
