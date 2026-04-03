package http

import (
	"net/http"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	useCase *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: uc}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Amount  int64  `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := &domain.Payment{
		ID:      uuid.New().String(),
		OrderID: req.OrderID,
		Amount:  req.Amount,
	}

	tid, status, err := h.useCase.ProcessPayment(c.Request.Context(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": tid,
		"status":         status,
	})
}
