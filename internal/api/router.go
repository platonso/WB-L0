package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/order/{order_uid}", h.GetOrder)
	r.Post("/order", h.CreateOrder)

	return r
}
