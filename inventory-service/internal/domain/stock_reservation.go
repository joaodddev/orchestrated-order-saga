package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusReserved ReservationStatus = "RESERVED"
	ReservationStatusFailed   ReservationStatus = "FAILED"
)

type StockReservation struct {
	OrderID    string
	CustomerID string
	Status     ReservationStatus
	UpdatedAt  time.Time
}

func Reserve(orderID, customerID string) *StockReservation {
	status := ReservationStatusReserved
	if orderID[len(orderID)-1] == '0' {
		status = ReservationStatusFailed
	}

	return &StockReservation{
		OrderID:    orderID,
		CustomerID: customerID,
		Status:     status,
		UpdatedAt:  time.Now().UTC(),
	}
}

func (s *StockReservation) Reserved() bool {
	return s.Status == ReservationStatusReserved
}
