package service

import (
	"context"
	"github.com/platonso/order-viewer/internal/domain"
	"github.com/platonso/order-viewer/internal/repository"
)

type OrderService struct {
	repo repository.Repository
}

func NewOrderService(repo repository.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) SaveOrder(ctx context.Context, order *domain.Order) error {
	// TODO
	return s.repo.Save(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	//TODO
	return s.repo.FindByID(ctx, orderUID)
}
