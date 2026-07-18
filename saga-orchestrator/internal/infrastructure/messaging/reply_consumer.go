package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/usecase"
)

type replyEnvelope struct {
	ReplyType string `json:"replyType"`
	SagaID    string `json:"sagaId"`
	Success   bool   `json:"success"`
	Reason    string `json:"reason"`
}

// ReplyConsumer escuta MÚLTIPLOS tópicos .reply simultaneamente — diferente
// dos consumers da saga coreografada, que cada um só escutava UM tópico.
// Isso é o reflexo direto da centralização: um único componente processa
// todas as respostas da saga, não importa de qual serviço vieram.
type ReplyConsumer struct {
	reader  *kafka.Reader
	useCase *usecase.AdvanceSaga
}

func NewReplyConsumer(brokers []string, groupID string, topics []string, useCase *usecase.AdvanceSaga) []*ReplyConsumer {
	consumers := make([]*ReplyConsumer, len(topics))
	for i, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		})
		consumers[i] = &ReplyConsumer{reader: reader, useCase: useCase}
	}
	return consumers
}

func (c *ReplyConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("reply consumer: failed to read message: %v", err)
			continue
		}

		var envelope replyEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("reply consumer: failed to unmarshal message: %v", err)
			continue
		}

		sagaID, err := uuid.Parse(envelope.SagaID)
		if err != nil {
			log.Printf("reply consumer: invalid saga id: %v", err)
			continue
		}

		if err := c.useCase.Execute(ctx, sagaID, envelope.ReplyType, envelope.Success, envelope.Reason); err != nil {
			log.Printf("reply consumer: failed to advance saga: %v", err)
			continue
		}

		log.Printf("[%s] saga %s advanced", envelope.ReplyType, envelope.SagaID)
	}
}

func (c *ReplyConsumer) Close() error {
	return c.reader.Close()
}
