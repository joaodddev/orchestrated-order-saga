package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Upsert(ctx context.Context, orderID, customerID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, status, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = VALUES(updated_at)`,
		orderID, customerID, status, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert order: %w", err)
	}
	return nil
}
