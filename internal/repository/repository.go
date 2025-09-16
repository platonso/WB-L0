package repository

import (
	"context"
	"github.com/platonso/order-viewer/internal/domain"
)

type DBRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, orderUID string) (*domain.Order, error)
	Close()
}

type CacheRepository interface {
	Save(order *domain.Order)
	FindByID(orderUID string) (*domain.Order, bool)
}
