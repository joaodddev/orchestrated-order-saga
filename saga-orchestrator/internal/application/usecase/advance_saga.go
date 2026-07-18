package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/port/output"
	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/domain"
)

type AdvanceSaga struct {
	repository output.SagaRepository
}

func NewAdvanceSaga(repository output.SagaRepository) *AdvanceSaga {
	return &AdvanceSaga{repository: repository}
}

func (uc *AdvanceSaga) Execute(ctx context.Context, sagaID uuid.UUID, replyType string, success bool, reason string) error {
	saga, err := uc.repository.FindByID(ctx, sagaID)
	if err != nil {
		return fmt.Errorf("saga not found: %w", err)
	}

	var nextCommand *domain.Command

	switch replyType {
	case "payment.reserve.reply":
		nextCommand = uc.handlePaymentReserveReply(saga, success)
	case "inventory.reserve.reply":
		nextCommand = uc.handleInventoryReserveReply(saga, success)
	case "payment.refund.reply":
		nextCommand = uc.handlePaymentRefundReply(saga)
	case "order.confirm.reply":
		saga.MarkConfirmed()
	case "order.cancel.reply":
		saga.MarkCancelled()
	default:
		return fmt.Errorf("unknown reply type: %s", replyType)
	}

	var outboxEvent *output.OutboxEvent
	if nextCommand != nil {
		payload, err := json.Marshal(nextCommand)
		if err != nil {
			return err
		}
		outboxEvent = &output.OutboxEvent{
			ID:          uuid.New(),
			AggregateID: saga.ID,
			EventType:   nextCommand.CommandType,
			Payload:     payload,
			CreatedAt:   time.Now().UTC(),
		}
	}

	return uc.repository.UpdateWithOutboxEvent(ctx, saga, outboxEvent)
}

func (uc *AdvanceSaga) handlePaymentReserveReply(saga *domain.Saga, success bool) *domain.Command {
	if !success {
		saga.MarkPaymentFailed()
		saga.MarkCancelled() // nada a compensar ainda, payment foi o primeiro passo
		return nil
	}

	saga.MarkPaymentReserved()
	return &domain.Command{
		CommandType: "inventory.reserve.command",
		Version:     1,
		SagaID:      saga.ID.String(),
		IssuedAt:    time.Now().UTC(),
		Payload: domain.ReserveStockCommandPayload{
			OrderID:    saga.OrderID.String(),
			CustomerID: saga.CustomerID.String(),
		},
	}
}

func (uc *AdvanceSaga) handleInventoryReserveReply(saga *domain.Saga, success bool) *domain.Command {
	if !success {
		saga.MarkStockFailed()
		saga.MarkCompensating()
		// Estoque falhou depois do pagamento já reservado — precisa compensar.
		return &domain.Command{
			CommandType: "payment.refund.command",
			Version:     1,
			SagaID:      saga.ID.String(),
			IssuedAt:    time.Now().UTC(),
			Payload: domain.RefundPaymentCommandPayload{
				OrderID: saga.OrderID.String(),
				Reason:  "insufficient stock",
			},
		}
	}

	saga.MarkStockReserved()
	return &domain.Command{
		CommandType: "order.confirm.command",
		Version:     1,
		SagaID:      saga.ID.String(),
		IssuedAt:    time.Now().UTC(),
		Payload: domain.ConfirmOrderCommandPayload{
			OrderID:    saga.OrderID.String(),
			CustomerID: saga.CustomerID.String(),
		},
	}
}

func (uc *AdvanceSaga) handlePaymentRefundReply(saga *domain.Saga) *domain.Command {
	return &domain.Command{
		CommandType: "order.cancel.command",
		Version:     1,
		SagaID:      saga.ID.String(),
		IssuedAt:    time.Now().UTC(),
		Payload: domain.CancelOrderCommandPayload{
			OrderID: saga.OrderID.String(),
			Reason:  "stock reservation failed",
		},
	}
}
