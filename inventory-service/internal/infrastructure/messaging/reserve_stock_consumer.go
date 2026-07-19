package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/application/usecase"
)

type commandEnvelope struct {
	CommandType string          `json:"commandType"`
	SagaID      string          `json:"sagaId"`
	Payload     json.RawMessage `json:"payload"`
}

type stockPayload struct {
	OrderID    string `json:"orderId"`
	CustomerID string `json:"customerId"`
}

type replyMessage struct {
	ReplyType  string    `json:"replyType"`
	SagaID     string    `json:"sagaId"`
	Success    bool      `json:"success"`
	Reason     string    `json:"reason,omitempty"`
	OccurredAt time.Time `json:"occurredAt"`
}

type ReserveStockConsumer struct {
	reader   *kafka.Reader
	producer *Producer
	useCase  *usecase.ReserveStock
}

func NewReserveStockConsumer(brokers []string, groupID string, producer *Producer, useCase *usecase.ReserveStock) *ReserveStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   "inventory.reserve.command",
	})
	return &ReserveStockConsumer{reader: reader, producer: producer, useCase: useCase}
}

func (c *ReserveStockConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("reserve stock consumer: failed to read message: %v", err)
			continue
		}

		var envelope commandEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("reserve stock consumer: failed to unmarshal command: %v", err)
			continue
		}

		var payload stockPayload
		if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
			log.Printf("reserve stock consumer: failed to unmarshal payload: %v", err)
			continue
		}

		reservation, err := c.useCase.Execute(ctx, payload.OrderID, payload.CustomerID)

		reply := replyMessage{
			ReplyType:  "inventory.reserve.reply",
			SagaID:     envelope.SagaID,
			OccurredAt: time.Now().UTC(),
		}
		if err != nil {
			reply.Success = false
			reply.Reason = err.Error()
		} else {
			reply.Success = reservation.Reserved()
			if !reply.Success {
				reply.Reason = "insufficient stock"
			}
		}

		replyPayload, _ := json.Marshal(reply)
		if err := c.producer.Publish(ctx, "inventory.reserve.reply", envelope.SagaID, replyPayload); err != nil {
			log.Printf("reserve stock consumer: failed to publish reply: %v", err)
			continue
		}

		log.Printf("[inventory.reserve.command] order %s -> success=%v", payload.OrderID, reply.Success)
	}
}

func (c *ReserveStockConsumer) Close() error {
	return c.reader.Close()
}
