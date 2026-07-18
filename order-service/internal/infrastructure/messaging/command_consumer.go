package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/application/usecase"
	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/domain"
)

type commandEnvelope struct {
	CommandType string          `json:"commandType"`
	SagaID      string          `json:"sagaId"`
	Payload     json.RawMessage `json:"payload"`
}

type orderPayload struct {
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

type CommandConsumer struct {
	reader   *kafka.Reader
	producer *Producer
	useCase  *usecase.ProcessOrderCommand
	status   domain.OrderStatus
	replyTo  string
}

func NewCommandConsumer(brokers []string, groupID, commandTopic, replyTopic string, status domain.OrderStatus, producer *Producer, useCase *usecase.ProcessOrderCommand) *CommandConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   commandTopic,
	})
	return &CommandConsumer{reader: reader, producer: producer, useCase: useCase, status: status, replyTo: replyTopic}
}

func (c *CommandConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("command consumer: failed to read message: %v", err)
			continue
		}

		var envelope commandEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("command consumer: failed to unmarshal command: %v", err)
			continue
		}

		var payload orderPayload
		if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
			log.Printf("command consumer: failed to unmarshal payload: %v", err)
			continue
		}

		err = c.useCase.Execute(ctx, payload.OrderID, payload.CustomerID, c.status)

		reply := replyMessage{
			ReplyType:  c.replyTo,
			SagaID:     envelope.SagaID,
			Success:    err == nil,
			OccurredAt: time.Now().UTC(),
		}
		if err != nil {
			reply.Reason = err.Error()
		}

		replyPayload, _ := json.Marshal(reply)
		if err := c.producer.Publish(ctx, c.replyTo, envelope.SagaID, replyPayload); err != nil {
			log.Printf("command consumer: failed to publish reply: %v", err)
			continue
		}

		log.Printf("[%s] order %s -> %s", envelope.CommandType, payload.OrderID, c.status)
	}
}

func (c *CommandConsumer) Close() error {
	return c.reader.Close()
}
