package main

import (
	"context"
	"log"

	"github.com/platonso/order-viewer/internal/app"
	"github.com/platonso/order-viewer/internal/config"
	"github.com/platonso/order-viewer/migrations"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if err := migrations.Run(ctx, cfg.GetConnStr()); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	defer application.DB.Close()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("Server run error: %v", err)
	}
}
