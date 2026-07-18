package domain

import (
	"time"

	"github.com/google/uuid"
)

type SagaStatus string

const (
	SagaStatusStarted         SagaStatus = "STARTED"
	SagaStatusPaymentReserved SagaStatus = "PAYMENT_RESERVED"
	SagaStatusPaymentFailed   SagaStatus = "PAYMENT_FAILED"
	SagaStatusStockReserved   SagaStatus = "STOCK_RESERVED"
	SagaStatusStockFailed     SagaStatus = "STOCK_FAILED"
	SagaStatusCompensating    SagaStatus = "COMPENSATING"
	SagaStatusConfirmed       SagaStatus = "CONFIRMED"
	SagaStatusCancelled       SagaStatus = "CANCELLED"
)

type Saga struct {
	ID          uuid.UUID
	OrderID     uuid.UUID
	CustomerID  uuid.UUID
	TotalAmount float64
	Status      SagaStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSaga(orderID, customerID uuid.UUID, totalAmount float64) *Saga {
	now := time.Now().UTC()
	return &Saga{
		ID:          uuid.New(),
		OrderID:     orderID,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Status:      SagaStatusStarted,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *Saga) MarkPaymentReserved() {
	s.Status = SagaStatusPaymentReserved
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkPaymentFailed() {
	s.Status = SagaStatusPaymentFailed
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkStockReserved() {
	s.Status = SagaStatusStockReserved
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkStockFailed() {
	s.Status = SagaStatusStockFailed
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkCompensating() {
	s.Status = SagaStatusCompensating
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkConfirmed() {
	s.Status = SagaStatusConfirmed
	s.UpdatedAt = time.Now().UTC()
}

func (s *Saga) MarkCancelled() {
	s.Status = SagaStatusCancelled
	s.UpdatedAt = time.Now().UTC()
}
