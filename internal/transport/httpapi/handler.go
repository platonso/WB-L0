package httpapi

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/platonso/order-viewer/internal/domain"
	"github.com/platonso/order-viewer/internal/service"
	"net/http"
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
		http.Error(w, "failed to save order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	order, err := h.orderService.GetOrder(r.Context(), orderUID)
	if err != nil {
		http.Error(w, domain.ErrOrderNotFound.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
