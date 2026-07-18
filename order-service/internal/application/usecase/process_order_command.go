package usecase

import (
	"context"

	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/orchestrated-order-saga/order-service/internal/domain"
)

type ProcessOrderCommand struct {
	repository output.OrderRepository
}

func NewProcessOrderCommand(repository output.OrderRepository) *ProcessOrderCommand {
	return &ProcessOrderCommand{repository: repository}
}

func (uc *ProcessOrderCommand) Execute(ctx context.Context, orderID, customerID string, status domain.OrderStatus) error {
	return uc.repository.Upsert(ctx, orderID, customerID, string(status))
}
