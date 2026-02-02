package product

import (
	"net/http"

	"github.com/gin-gonic/gin"

	productapp "flash-sale-order-system/internal/application/product"
)

type Handler struct {
	createHandler *productapp.CreateProductHandler
}

func NewHandler(createHandler *productapp.CreateProductHandler) *Handler {
	return &Handler{
		createHandler: createHandler,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/product")
	{
		products.POST("", h.Create)
		// products.GET("/:id", h.GetByID)
		// products.PUT("/:id", h.Update)
		// products.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := productapp.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Quantity:    req.Quantity,
		Prices:      req.Prices,
		PriceFrom:   req.PriceFrom,
		PriceUntil:  req.PriceUntil,
	}

	productID, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateProductResponse{ID: productID})
}
