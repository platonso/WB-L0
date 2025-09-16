package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/platonso/order-viewer/internal/domain"
	"github.com/platonso/order-viewer/internal/service"
	"net/http"
	"time"
)

type Handler struct {
	orderService *service.OrderService
}

func NewHandler(orderService *service.OrderService) *Handler {
	return &Handler{orderService: orderService}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.orderService.SaveOrder(r.Context(), &order); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrOrderAlreadyExists):
			http.Error(w, "order already exists", http.StatusConflict)
		default:
			http.Error(w, "failed to save order", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // замер начала обработки

	orderUID := chi.URLParam(r, "order_uid")

	order, fromCache, err := h.orderService.GetOrder(r.Context(), orderUID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrOrderNotFound):
			http.Error(w, "order not found", http.StatusNotFound)
		default:
			http.Error(w, "failed to fetch order", http.StatusInternalServerError)
		}
		return
	}

	if fromCache {
		w.Header().Set("X-Data-Source", "cache")
	} else {
		w.Header().Set("X-Data-Source", "database")
	}

	// Время ответа сервера
	duration := time.Since(start)
	// Преобразуем в миллисекунды как число с точкой
	ms := float64(duration.Nanoseconds()) / float64(time.Millisecond)
	w.Header().Set("X-Response-Time", fmt.Sprintf("%.3f", ms))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
