package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/platonso/order-viewer/internal/domain"
	"github.com/segmentio/kafka-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	brokersCSV := flag.String("brokers", "localhost:19092", "Kafka brokers (comma-separated), e.g. localhost:19092 or kafka:9092")
	topic := flag.String("topic", "orders", "Kafka topic")
	numOrders := flag.Int("orders", 10, "Number of orders to send")
	interval := flag.Duration("interval", 300*time.Millisecond, "Interval between messages")
	flag.Parse()

	brokers := splitCSV(*brokersCSV)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    *topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	log.Printf("Sending %d orders to %v topic %q...", *numOrders, brokers, *topic)

	for i := 0; i < *numOrders; i++ {
		order := generateOrder()
		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("failed to marshal order: %v", err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: data,
		})
		if err != nil {
			log.Printf("failed to send order: %v", err)
			continue
		}
		time.Sleep(*interval)
	}
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{s}
	}
	return out
}

func generateOrder() domain.Order {
	uid := fmt.Sprintf("test-%d", rand.Intn(1_000_000))

	items := make([]domain.Item, rand.Intn(5)+1) // 1..5
	for i := range items {
		items[i] = domain.Item{
			ChrtID:      rand.Intn(100000),
			TrackNumber: fmt.Sprintf("TRACK-%04d", rand.Intn(10000)),
			Price:       rand.Intn(1000) + 1,
			RID:         fmt.Sprintf("rid-%d", rand.Intn(100000)),
			Name:        fmt.Sprintf("Product-%d", i+1),
			Sale:        rand.Intn(50),
			Size:        "M",
			TotalPrice:  rand.Intn(2000) + 1,
			NmID:        rand.Intn(10000),
			Brand:       "BrandX",
			Status:      200,
		}
	}

	return domain.Order{
		OrderUID:    uid,
		TrackNumber: fmt.Sprintf("TRACK-%04d", rand.Intn(10000)),
		Entry:       "WBIL",
		Locale:      "en",
		Delivery: domain.Delivery{
			Name:    "Platon Solozobov",
			Phone:   "+123456789",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Lenina 1",
			Region:  "Moscow",
			Email:   "test@example.com",
		},
		Payment: domain.Payment{
			Transaction:  uid,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 1,
			PaymentDt:    int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   rand.Intn(3000) + 1,
			CustomFee:    0,
		},
		Items:           items,
		CustomerID:      "test-customer",
		DeliveryService: "meest",
		Shardkey:        "9",
		SmID:            78,
		DateCreated:     time.Now().UTC(),
		OofShard:        "1",
	}
}
