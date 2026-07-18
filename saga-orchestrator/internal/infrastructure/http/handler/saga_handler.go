package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/application/usecase"
)

type SagaHandler struct {
	startSaga *usecase.StartSaga
}

func NewSagaHandler(startSaga *usecase.StartSaga) *SagaHandler {
	return &SagaHandler{startSaga: startSaga}
}

type startSagaRequest struct {
	OrderID     string  `json:"orderId" binding:"required"`
	CustomerID  string  `json:"customerId" binding:"required"`
	TotalAmount float64 `json:"totalAmount" binding:"required,gt=0"`
}

func (h *SagaHandler) Start(c *gin.Context) {
	var req startSagaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orderId"})
		return
	}
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customerId"})
		return
	}

	saga, err := h.startSaga.Execute(c.Request.Context(), usecase.StartSagaInput{
		OrderID:     orderID,
		CustomerID:  customerID,
		TotalAmount: req.TotalAmount,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"sagaId": saga.ID,
		"status": saga.Status,
	})
}
