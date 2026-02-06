package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"flash-sale-order-system/internal/application/product/query"
)

type QueryHandler struct {
	queryHandler *query.ProductQueryHandler
}

func NewQueryHandler(
	queryHandler *query.ProductQueryHandler,
) *QueryHandler {
	return &QueryHandler{
		queryHandler: queryHandler,
	}
}

func (h *QueryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := h.queryHandler.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
