package http

import (
	"github.com/gin-gonic/gin"

	"github.com/joaodddev/orchestrated-order-saga/saga-orchestrator/internal/infrastructure/http/handler"
)

func NewRouter(sagaHandler *handler.SagaHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/sagas", sagaHandler.Start)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	return router
}
