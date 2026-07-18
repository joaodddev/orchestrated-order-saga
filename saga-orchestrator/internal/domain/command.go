package domain

import "time"

type Command struct {
	CommandType string      `json:"commandType"`
	Version     int         `json:"version"`
	SagaID      string      `json:"sagaId"`
	IssuedAt    time.Time   `json:"issuedAt"`
	Payload     interface{} `json:"payload"`
}

type ReservePaymentCommandPayload struct {
	OrderID     string  `json:"orderId"`
	CustomerID  string  `json:"customerId"`
	TotalAmount float64 `json:"totalAmount"`
}

type ReserveStockCommandPayload struct {
	OrderID    string `json:"orderId"`
	CustomerID string `json:"customerId"`
}

type RefundPaymentCommandPayload struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
}

type ConfirmOrderCommandPayload struct {
	OrderID    string `json:"orderId"`
	CustomerID string `json:"customerId"`
}

type CancelOrderCommandPayload struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
}
