package order

import (
	"fmt"
	"net/http"

	apporder "flash-sale-order-system/internal/application/order"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	placeOrder *apporder.PlaceOrderHandler
}

func NewHandler(placeOrder *apporder.PlaceOrderHandler) *Handler {
	return &Handler{placeOrder: placeOrder}
}

type placeOrderRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	ProductID int64   `json:"product_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required"`
}

func (h *Handler) PlaceOrder(c *gin.Context) {
	var req placeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := apporder.PlaceOrderCommand{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}

	orderID, err := h.placeOrder.Handle(c.Request.Context(), cmd)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == fmt.Sprintf("insufficient stock for product %d", req.ProductID) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"order_id": orderID, "status": "queued"})
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/orders", h.PlaceOrder)
}
