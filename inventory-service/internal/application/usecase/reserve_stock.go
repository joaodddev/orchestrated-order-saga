package usecase

import (
	"context"

	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/application/port/output"
	"github.com/joaodddev/orchestrated-order-saga/inventory-service/internal/domain"
)

type ReserveStock struct {
	repository output.StockRepository
}

func NewReserveStock(repository output.StockRepository) *ReserveStock {
	return &ReserveStock{repository: repository}
}

func (uc *ReserveStock) Execute(ctx context.Context, orderID, customerID string) (*domain.StockReservation, error) {
	reservation := domain.Reserve(orderID, customerID)
	if err := uc.repository.Upsert(ctx, orderID, customerID, string(reservation.Status)); err != nil {
		return nil, err
	}
	return reservation, nil
}
