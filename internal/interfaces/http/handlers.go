package http

import (
	"flash-sale-order-system/internal/interfaces/http/product"
)

type Handlers struct {
	Product *product.Handler
}
