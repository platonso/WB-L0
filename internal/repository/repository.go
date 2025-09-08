package repository

import (
	"context"
	"github.com/platonso/order-viewer/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, orderUID string) (*domain.Order, error)
	Close()
}
