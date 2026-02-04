package query

import (
	"context"
)

type ProductQueryHandler struct {
	queryService ProductQueryService
}

type ProductQueryService interface {
	GetByID(ctx context.Context, id int64) (*ProductDTO, error)
	GetWithCurrentPrice(ctx context.Context, id int64) (*ProductWithPriceDTO, error)
}

func NewProductQueryHandler(queryService ProductQueryService) *ProductQueryHandler {
	return &ProductQueryHandler{
		queryService: queryService,
	}
}

func (h *ProductQueryHandler) GetByID(ctx context.Context, id int64) (*ProductDTO, error) {
	return h.queryService.GetByID(ctx, id)
}

func (h *ProductQueryHandler) GetWithCurrentPrice(ctx context.Context, id int64) (*ProductWithPriceDTO, error) {
	return h.queryService.GetWithCurrentPrice(ctx, id)
}
