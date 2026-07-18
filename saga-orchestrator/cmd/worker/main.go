package main

import (
	"context"
	"log"
	"os"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/usecase"
	httpinfra "github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/infrastructure/http"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/infrastructure/http/handler"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/infrastructure/messaging"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/infrastructure/persistence/mysql"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/outbox"
)

var replyTopics = []string{
	"payment.reserve.reply",
	"inventory.reserve.reply",
	"payment.refund.reply",
	"order.confirm.reply",
	"order.cancel.reply",
}

func main() {
	ctx := context.Background()

	db, err := mysql.NewConnection(mysql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3308"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_NAME", "orchestrator_db"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	kafkaBrokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}

	producer := messaging.NewProducer(kafkaBrokers)
	defer producer.Close()

	outboxRepository := mysql.NewOutboxRepository(db)
	relay := outbox.NewRelay(outboxRepository, producer)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go relay.Start(ctx)

	sagaRepository := mysql.NewSagaRepository(db)

	startSagaUseCase := usecase.NewStartSaga(sagaRepository)
	sagaHandler := handler.NewSagaHandler(startSagaUseCase)

	advanceSagaUseCase := usecase.NewAdvanceSaga(sagaRepository)
	replyConsumers := messaging.NewReplyConsumer(kafkaBrokers, "saga-orchestrator-group", replyTopics, advanceSagaUseCase)
	for _, consumer := range replyConsumers {
		c := consumer
		defer c.Close()
		go c.Start(ctx)
	}

	router := httpinfra.NewRouter(sagaHandler)

	port := getEnv("PORT", "8090")
	log.Printf("saga-orchestrator listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
