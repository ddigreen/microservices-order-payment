package http

import (
	"net/http"
	"strconv"
	"time"

	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	useCase *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: uc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		ItemName   string `json:"item_name" binding:"required"`
		Amount     int64  `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: req.CustomerID,
		ItemName:   req.ItemName,
		Amount:     req.Amount,
		CreatedAt:  time.Now(),
	}

	if err := h.useCase.CreateOrder(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.useCase.CancelOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Cancelled"})
}

func (h *OrderHandler) GetRecent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "limit must be a valid integer"})
		return
	}

	orders, err := h.useCase.GetRecentOrders(c.Request.Context(), limit)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}
