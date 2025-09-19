package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/platonso/order-viewer/internal/api"
	"github.com/platonso/order-viewer/internal/config"
	"github.com/platonso/order-viewer/internal/kafka"
	"github.com/platonso/order-viewer/internal/repository"
	"github.com/platonso/order-viewer/internal/service"
)

type Application struct {
	Config *config.Config
	DB     repository.DBRepository
	Cache  repository.CacheRepository
}

func NewApp(ctx context.Context, cfg *config.Config) (*Application, error) {
	postgresRepo, err := repository.NewPostgresRepo(ctx, cfg.GetConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to create order repository: %w", err)
	}

	cacheRepo := repository.NewCacheRepo()

	return &Application{
		Config: cfg,
		DB:     postgresRepo,
		Cache:  cacheRepo,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	orderService := service.NewOrderService(app.DB, app.Cache)

	handler := api.NewHandler(orderService)
	router := api.NewRouter(handler)

	if err := kafka.StartConsumer(ctx, app.Config, orderService); err != nil {
		log.Printf("failed to start kafka consumer: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: router,
	}

	log.Printf("Server is running on port %s", app.Config.Port)

	return srv.ListenAndServe()
}
