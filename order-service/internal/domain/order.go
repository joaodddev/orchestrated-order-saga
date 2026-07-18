package domain

import "time"

type OrderStatus string

const (
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID         string
	CustomerID string
	Status     OrderStatus
	UpdatedAt  time.Time
}
