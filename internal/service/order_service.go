package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/platonso/order-viewer/internal/domain"
	"github.com/platonso/order-viewer/internal/repository"
	"log"
	"regexp"
	"strings"
	"time"
)

type OrderService struct {
	dbRepo    repository.DBRepository
	cacheRepo repository.CacheRepository
}

func NewOrderService(dbRepo repository.DBRepository, cacheRepo repository.CacheRepository) *OrderService {
	return &OrderService{
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *OrderService) SaveOrder(ctx context.Context, order *domain.Order) error {

	// Валидация всех полей заказа
	if err := s.validateOrder(order); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrValidation, err)
	}

	// Сохранение нового заказа в бд
	if err := s.dbRepo.Save(ctx, order); err != nil {
		log.Printf("failed to save order in db: %v", err)
		return err
	}
	// Добавление в кэш
	s.cacheRepo.Save(order)

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*domain.Order, bool, error) {

	// Валидация uid заказа
	if orderUID == "" {
		return nil, false, fmt.Errorf("%w: order id is required", domain.ErrValidation)
	}

	if len(orderUID) > 36 {
		return nil, false, fmt.Errorf("%w: order id is too long", domain.ErrValidation)
	}

	validID := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validID.MatchString(orderUID) {
		return nil, false, fmt.Errorf("%w: order id contains invalid characters", domain.ErrValidation)
	}

	// Попытка достать заказ из кэша
	order, ok := s.cacheRepo.FindByID(orderUID)
	if ok {
		return order, true, nil
	}

	// При отсутствии заказа в кэше, поиск его в бд
	order, err := s.dbRepo.FindByID(ctx, orderUID)
	if err != nil {
		return nil, false, err
	}

	// Добавление заказа в кэш, если он нашёлся в бд
	if order != nil {
		s.cacheRepo.Save(order)
		return order, false, nil
	}

	return nil, false, domain.ErrOrderNotFound
}

func (s *OrderService) validateOrder(order *domain.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	if order.OrderUID == "" {
		return errors.New("order_uid is required")
	}
	if len(order.OrderUID) > 36 {
		return errors.New("order_uid is too long (max 36 characters)")
	}
	if strings.Contains(order.OrderUID, " ") {
		return errors.New("order_uid cannot contain spaces")
	}

	if order.TrackNumber == "" {
		return errors.New("track_number is required")
	}
	if len(order.TrackNumber) > 36 {
		return errors.New("track_number is too long (max 36 characters)")
	}

	if order.DateCreated.IsZero() {
		return errors.New("date_created is required")
	}
	if order.DateCreated.After(time.Now().Add(time.Hour)) {
		return errors.New("date_created cannot be in the future")
	}

	if err := s.validateDelivery(&order.Delivery); err != nil {
		return fmt.Errorf("delivery validation failed: %w", err)
	}

	if err := s.validatePayment(&order.Payment); err != nil {
		return fmt.Errorf("payment validation failed: %w", err)
	}

	if len(order.Items) == 0 {
		return errors.New("at least one item is required")
	}
	for i, item := range order.Items {
		if err := s.validateItem(&item); err != nil {
			return fmt.Errorf("item[%d] validation failed: %w", i, err)
		}
	}

	return nil
}

func (s *OrderService) validateDelivery(delivery *domain.Delivery) error {
	if delivery == nil {
		return errors.New("delivery is nil")
	}

	if delivery.Name == "" {
		return errors.New("delivery name is required")
	}
	if len(delivery.Name) > 100 {
		return errors.New("delivery name is too long (max 100 characters)")
	}

	if delivery.Phone == "" {
		return errors.New("delivery phone is required")
	}

	if delivery.Email == "" {
		return errors.New("delivery email is required")
	}
	if len(delivery.Email) > 254 {
		return errors.New("delivery email is too long (max 254 characters)")
	}
	if !strings.Contains(delivery.Email, "@") || !strings.Contains(delivery.Email, ".") {
		return errors.New("delivery email format is invalid")
	}

	return nil
}

func (s *OrderService) validatePayment(payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment is nil")
	}

	if payment.Transaction == "" {
		return errors.New("payment transaction is required")
	}
	if len(payment.Transaction) > 36 {
		return errors.New("payment transaction is too long (max 36 characters)")
	}

	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than 0")
	}

	if payment.Currency == "" {
		return errors.New("payment currency is required")
	}
	if len(payment.Currency) > 3 {
		return errors.New("payment currency code is invalid (max 3 characters)")
	}

	return nil
}

func (s *OrderService) validateItem(item *domain.Item) error {
	if item == nil {
		return errors.New("item is nil")
	}

	if item.Name == "" {
		return errors.New("item name is required")
	}
	if len(item.Name) > 200 {
		return errors.New("item name is too long (max 200 characters)")
	}

	if item.Price <= 0 {
		return errors.New("item price must be greater than 0")
	}
	if item.Price > 100000000 {
		return errors.New("item price is too large")
	}

	if item.TotalPrice < 0 {
		return errors.New("item total_price cannot be negative")
	}

	if item.ChrtID <= 0 {
		return errors.New("item chrt_id must be greater than 0")
	}

	return nil
}
