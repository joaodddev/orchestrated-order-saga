package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/port/output"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/domain"
)

type SagaRepository struct {
	db *sql.DB
}

func NewSagaRepository(db *sql.DB) *SagaRepository {
	return &SagaRepository{db: db}
}

func (r *SagaRepository) SaveWithOutboxEvent(ctx context.Context, saga *domain.Saga, event output.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO saga_state (id, order_id, customer_id, total_amount, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		saga.ID.String(), saga.OrderID.String(), saga.CustomerID.String(),
		saga.TotalAmount, saga.Status, saga.CreatedAt, saga.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert saga: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events (id, aggregate_id, event_type, payload, published, created_at)
		 VALUES (?, ?, ?, ?, FALSE, ?)`,
		event.ID.String(), event.AggregateID.String(), event.EventType, event.Payload, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert outbox event: %w", err)
	}

	return tx.Commit()
}

func (r *SagaRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Saga, error) {
	return r.scanOne(ctx, `SELECT id, order_id, customer_id, total_amount, status, created_at, updated_at
		FROM saga_state WHERE id = ?`, id.String())
}

func (r *SagaRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Saga, error) {
	return r.scanOne(ctx, `SELECT id, order_id, customer_id, total_amount, status, created_at, updated_at
		FROM saga_state WHERE order_id = ?`, orderID.String())
}

func (r *SagaRepository) scanOne(ctx context.Context, query string, arg string) (*domain.Saga, error) {
	row := r.db.QueryRowContext(ctx, query, arg)

	var saga domain.Saga
	var id, orderID, customerID string
	if err := row.Scan(&id, &orderID, &customerID, &saga.TotalAmount, &saga.Status, &saga.CreatedAt, &saga.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("saga not found")
		}
		return nil, fmt.Errorf("failed to query saga: %w", err)
	}

	saga.ID = uuid.MustParse(id)
	saga.OrderID = uuid.MustParse(orderID)
	saga.CustomerID = uuid.MustParse(customerID)

	return &saga, nil
}

func (r *SagaRepository) UpdateWithOutboxEvent(ctx context.Context, saga *domain.Saga, event *output.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE saga_state SET status = ?, updated_at = ? WHERE id = ?`,
		saga.Status, saga.UpdatedAt, saga.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update saga: %w", err)
	}

	if event != nil {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO outbox_events (id, aggregate_id, event_type, payload, published, created_at)
			 VALUES (?, ?, ?, ?, FALSE, ?)`,
			event.ID.String(), event.AggregateID.String(), event.EventType, event.Payload, event.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert outbox event: %w", err)
		}
	}

	return tx.Commit()
}
