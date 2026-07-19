package output

import "context"

type StockRepository interface {
	Upsert(ctx context.Context, orderID, customerID, status string) error
}
