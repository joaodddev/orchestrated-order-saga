package output

import "context"

type OrderRepository interface {
	Upsert(ctx context.Context, orderID, customerID, status string) error
}
