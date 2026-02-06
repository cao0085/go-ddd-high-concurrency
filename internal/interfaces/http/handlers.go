package http

import (
	"flash-sale-order-system/internal/interfaces/http/product"
)

type Handlers struct {
	ProductCommand *product.CommandHandler
	ProductQuery   *product.QueryHandler
}
