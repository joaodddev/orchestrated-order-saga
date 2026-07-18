package main

import (
	"context"
	"log"
	"os"

	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/application/usecase"
	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/domain"
	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/infrastructure/messaging"
	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/infrastructure/persistence/mysql"
)

func main() {
	ctx := context.Background()

	db, err := mysql.NewConnection(mysql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3309"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_NAME", "order_db"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	kafkaBrokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	producer := messaging.NewProducer(kafkaBrokers)
	defer producer.Close()

	orderRepository := mysql.NewOrderRepository(db)
	processCommand := usecase.NewProcessOrderCommand(orderRepository)

	confirmConsumer := messaging.NewCommandConsumer(
		kafkaBrokers, "order-service-confirm-group",
		"order.confirm.command", "order.confirm.reply",
		domain.OrderStatusConfirmed, producer, processCommand,
	)
	defer confirmConsumer.Close()
	go confirmConsumer.Start(ctx)

	cancelConsumer := messaging.NewCommandConsumer(
		kafkaBrokers, "order-service-cancel-group",
		"order.cancel.command", "order.cancel.reply",
		domain.OrderStatusCancelled, producer, processCommand,
	)
	defer cancelConsumer.Close()

	log.Println("order-service worker started, consuming order commands")
	cancelConsumer.Start(ctx)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
