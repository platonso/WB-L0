package repository

import (
	"github.com/platonso/order-viewer/internal/domain"
	"sync"
)

type CacheRepo struct {
	orders map[string]*domain.Order
	mu     sync.RWMutex
}

func NewCacheRepo() *CacheRepo {
	return &CacheRepo{
		orders: make(map[string]*domain.Order),
	}
}

func (c *CacheRepo) Save(order *domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *CacheRepo) FindByID(orderUID string) (*domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[orderUID]
	return order, ok
}
