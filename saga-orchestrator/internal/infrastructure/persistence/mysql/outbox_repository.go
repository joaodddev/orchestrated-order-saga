package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/port/output"
)

type OutboxRepository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) FetchPending(ctx context.Context, limit int) ([]output.PendingOutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, aggregate_id, event_type, payload
		 FROM outbox_events
		 WHERE published = FALSE
		 ORDER BY created_at ASC
		 LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending outbox events: %w", err)
	}
	defer rows.Close()

	var events []output.PendingOutboxEvent
	for rows.Next() {
		var e output.PendingOutboxEvent
		var id, aggregateID string
		if err := rows.Scan(&id, &aggregateID, &e.EventType, &e.Payload); err != nil {
			return nil, fmt.Errorf("failed to scan outbox event: %w", err)
		}
		e.ID = uuid.MustParse(id)
		e.AggregateID = uuid.MustParse(aggregateID)
		events = append(events, e)
	}

	return events, rows.Err()
}

func (r *OutboxRepository) MarkAsPublished(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE outbox_events SET published = TRUE WHERE id = ?`, id.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to mark outbox event as published: %w", err)
	}
	return nil
}
