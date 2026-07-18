package outbox

import (
	"context"
	"log"
	"time"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/port/output"
)

type Publisher interface {
	Publish(ctx context.Context, topic, key string, payload []byte) error
}

type Relay struct {
	repository output.OutboxRepository
	publisher  Publisher
	interval   time.Duration
	batchSize  int
}

func NewRelay(repository output.OutboxRepository, publisher Publisher) *Relay {
	return &Relay{
		repository: repository,
		publisher:  publisher,
		interval:   2 * time.Second,
		batchSize:  20,
	}
}

func (r *Relay) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("outbox relay stopped")
			return
		case <-ticker.C:
			r.processBatch(ctx)
		}
	}
}

func (r *Relay) processBatch(ctx context.Context) {
	events, err := r.repository.FetchPending(ctx, r.batchSize)
	if err != nil {
		log.Printf("outbox relay: failed to fetch pending events: %v", err)
		return
	}

	for _, event := range events {
		err := r.publisher.Publish(ctx, event.EventType, event.AggregateID.String(), event.Payload)
		if err != nil {
			log.Printf("outbox relay: failed to publish event %s: %v", event.ID, err)
			continue
		}

		if err := r.repository.MarkAsPublished(ctx, event.ID); err != nil {
			log.Printf("outbox relay: failed to mark event %s as published: %v", event.ID, err)
		}
	}
}
