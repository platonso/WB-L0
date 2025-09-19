package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/platonso/order-viewer/internal/config"
	"github.com/platonso/order-viewer/internal/domain"
	"github.com/platonso/order-viewer/internal/service"

	"github.com/segmentio/kafka-go"
)

// StartConsumer запускает чтение сообщений из Kafka и сохраняет заказы через сервис.
func StartConsumer(ctx context.Context, cfg *config.Config, orderService *service.OrderService) error {
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: cfg.KafkaGroupID,
		Topic:   cfg.KafkaTopic,
	})

	go func() {
		defer func() {
			if err := reader.Close(); err != nil {
				log.Printf("kafka reader close error: %v", err)
			}
		}()

		log.Printf("Kafka consumer started. Brokers=%s Topic=%s Group=%s", cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)

		for {
			select {
			case <-ctx.Done():
				log.Printf("Kafka consumer stopped: %v", ctx.Err())
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						log.Printf("Kafka consumer exiting: %v", ctx.Err())
						return
					}
					log.Printf("kafka read error: %v", err)
					continue
				}

				var order domain.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("kafka invalid message, json error: %v", err)
					continue
				}

				ctxSave, cancel := context.WithTimeout(ctx, 5*time.Second)
				if err := orderService.SaveOrder(ctxSave, &order); err != nil {
					log.Printf("failed to save order from kafka: %v", err)
				}
				cancel()
			}
		}
	}()

	return nil
}
