package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) Upsert(ctx context.Context, orderID, customerID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stock_reservations (order_id, customer_id, status, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = VALUES(updated_at)`,
		orderID, customerID, status, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert stock reservation: %w", err)
	}
	return nil
}
