package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	// Находим абсолютный путь к папке web
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	webDir := filepath.Join(wd, "internal/web")

	// Frontend
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	// Endpoints
	r.Get("/order/{order_uid}", h.GetOrder)
	r.Post("/order", h.CreateOrder)

	return r
}
