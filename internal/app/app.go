package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/order-viewer/internal/api"
	"github.com/platonso/order-viewer/internal/config"
	"github.com/platonso/order-viewer/internal/repository"
	"github.com/platonso/order-viewer/internal/service"
	"log"
	"net/http"
)

type Application struct {
	Config    *config.Config
	DB        repository.DBRepository
	Cache     repository.CacheRepository
	Validator *validator.Validate
}

func NewApp(ctx context.Context, cfg *config.Config) (*Application, error) {
	postgresRepo, err := repository.NewPostgresRepo(ctx, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order repository: %w", err)
	}

	cacheRepo := repository.NewCacheRepo()

	validate := validator.New()

	return &Application{
		Config:    cfg,
		DB:        postgresRepo,
		Cache:     cacheRepo,
		Validator: validate,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	orderService := service.NewOrderService(app.DB, app.Cache)

	handler := api.NewHandler(orderService)
	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: router,
	}

	log.Printf("Server is running on port %s", app.Config.Port)

	return srv.ListenAndServe()
}
