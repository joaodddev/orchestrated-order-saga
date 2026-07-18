package output

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/domain"
)

type OutboxEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}

type SagaRepository interface {
	SaveWithOutboxEvent(ctx context.Context, saga *domain.Saga, event OutboxEvent) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Saga, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Saga, error)

	UpdateWithOutboxEvent(ctx context.Context, saga *domain.Saga, event *OutboxEvent) error
}
