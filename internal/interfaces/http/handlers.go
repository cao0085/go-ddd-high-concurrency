package http

import (
	"flash-sale-order-system/internal/interfaces/http/product"
)

type Handlers struct {
	Product *product.Handler
}

func NewHandlers(
	productHandler *product.Handler,
) *Handlers {
	return &Handlers{
		Product: productHandler,
	}
}
