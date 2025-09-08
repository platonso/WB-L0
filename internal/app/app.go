package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/order-viewer/internal/config"
	"github.com/platonso/order-viewer/internal/repository"
	"github.com/platonso/order-viewer/internal/service"
	"github.com/platonso/order-viewer/internal/transport/httpapi"
	"log"
	"net/http"
)

type Application struct {
	Config    *config.Config
	DB        repository.Repository
	Validator *validator.Validate
	//	KafkaConsumer *kafka.Consumer
}

func NewApp(ctx context.Context, cfg *config.Config) (*Application, error) {
	postgresRepo, err := repository.NewPostgresRepo(ctx, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order repository: %w", err)
	}

	validate := validator.New()

	return &Application{
		Config:    cfg,
		DB:        postgresRepo,
		Validator: validate,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	orderService := service.NewOrderService(a.DB)

	handler := httpapi.NewHandler(orderService)
	router := httpapi.NewRouter(handler)

	srv := &http.Server{
		Addr:    ":" + a.Config.Port,
		Handler: router,
	}

	log.Printf("Server is running on port %s", a.Config.Port)

	return srv.ListenAndServe()
}
