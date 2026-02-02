package product

import (
	"time"
)

type CreateProductRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	SKU         string             `json:"sku" binding:"required"`
	Quantity    int32              `json:"quantity" binding:"required,min=0"`
	Prices      map[string]float64 `json:"prices" binding:"required"`
	PriceFrom   time.Time          `json:"price_from" binding:"required"`
	PriceUntil  *time.Time         `json:"price_until"`
}
