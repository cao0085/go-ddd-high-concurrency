package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	productapp "flash-sale-order-system/internal/application/product"
	shareddomain "flash-sale-order-system/internal/shared/domain"
)

type ProductHandler struct {
	createProductHandler *productapp.CreateProductHandler
}

func NewProductHandler(createProductHandler *productapp.CreateProductHandler) *ProductHandler {
	return &ProductHandler{
		createProductHandler: createProductHandler,
	}
}

type CreateProductRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	SKU         string             `json:"sku" binding:"required"`
	Quantity    int32              `json:"quantity" binding:"required,min=0"`
	Prices      map[string]float64 `json:"prices" binding:"required"`
	PriceFrom   time.Time          `json:"price_from" binding:"required"`
	PriceUntil  *time.Time         `json:"price_until"`
}

type CreateProductResponse struct {
	ID int64 `json:"id"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert prices map
	prices := make(map[shareddomain.Currency]shareddomain.Money)
	for currencyStr, amount := range req.Prices {
		currency := shareddomain.Currency(currencyStr)
		money, err := shareddomain.NewMoney(amount, currency)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prices[currency] = money
	}

	cmd := productapp.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Quantity:    req.Quantity,
		Prices:      prices,
		PriceFrom:   req.PriceFrom,
		PriceUntil:  req.PriceUntil,
	}

	productID, err := h.createProductHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateProductResponse{ID: productID})
}
