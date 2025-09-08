package repository

import (
	"context"
	"github.com/platonso/order-viewer/internal/domain"
	"sync"
)

type CacheRepo struct {
	orders map[string]*domain.Order
	mu     sync.Mutex
}

func NewCache() *CacheRepo {
	return &CacheRepo{
		orders: make(map[string]*domain.Order),
	}
}

func (c *CacheRepo) Save(ctx context.Context, order *domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
	return nil
}

func (c *CacheRepo) FindByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	order, ok := c.orders[orderUID]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (c *CacheRepo) Close() {}
