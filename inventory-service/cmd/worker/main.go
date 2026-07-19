package main

import (
	"context"
	"log"
	"os"

	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/application/usecase"
	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/infrastructure/messaging"
	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/infrastructure/persistence/mysql"
)

func main() {
	ctx := context.Background()

	db, err := mysql.NewConnection(mysql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3310"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_NAME", "inventory_db"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	kafkaBrokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	producer := messaging.NewProducer(kafkaBrokers)
	defer producer.Close()

	stockRepository := mysql.NewStockRepository(db)
	reserveStock := usecase.NewReserveStock(stockRepository)

	consumer := messaging.NewReserveStockConsumer(kafkaBrokers, "inventory-service-group", producer, reserveStock)
	defer consumer.Close()

	log.Println("inventory-service worker started, consuming inventory.reserve.command")
	consumer.Start(ctx)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
