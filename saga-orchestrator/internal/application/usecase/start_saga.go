package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/port/output"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/domain"
)

type StartSagaInput struct {
	OrderID     uuid.UUID
	CustomerID  uuid.UUID
	TotalAmount float64
}

type StartSaga struct {
	repository output.SagaRepository
}

func NewStartSaga(repository output.SagaRepository) *StartSaga {
	return &StartSaga{repository: repository}
}

func (uc *StartSaga) Execute(ctx context.Context, in StartSagaInput) (*domain.Saga, error) {
	saga := domain.NewSaga(in.OrderID, in.CustomerID, in.TotalAmount)

	command := domain.Command{
		CommandType: "payment.reserve.command",
		Version:     1,
		SagaID:      saga.ID.String(),
		IssuedAt:    time.Now().UTC(),
		Payload: domain.ReservePaymentCommandPayload{
			OrderID:     saga.OrderID.String(),
			CustomerID:  saga.CustomerID.String(),
			TotalAmount: saga.TotalAmount,
		},
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	outboxEvent := output.OutboxEvent{
		ID:          uuid.New(),
		AggregateID: saga.ID,
		EventType:   command.CommandType,
		Payload:     payload,
		CreatedAt:   saga.CreatedAt,
	}

	if err := uc.repository.SaveWithOutboxEvent(ctx, saga, outboxEvent); err != nil {
		return nil, err
	}

	return saga, nil
}
