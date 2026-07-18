package output

import (
	"context"

	"github.com/google/uuid"
)

type PendingOutboxEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
}

type OutboxRepository interface {
	FetchPending(ctx context.Context, limit int) ([]PendingOutboxEvent, error)
	MarkAsPublished(ctx context.Context, id uuid.UUID) error
}
