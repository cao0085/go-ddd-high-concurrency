package query

import (
	"context"
)

type GetProductHandler struct {
	queryService ProductQueryService
}

func NewGetProductHandler(queryService ProductQueryService) *GetProductHandler {
	return &GetProductHandler{
		queryService: queryService,
	}
}

func (h *GetProductHandler) Handle(ctx context.Context, id int64) (*ProductDTO, error) {
	return h.queryService.GetByID(ctx, id)
}

type GetProductWithPriceHandler struct {
	queryService ProductQueryService
}

func NewGetProductWithPriceHandler(queryService ProductQueryService) *GetProductWithPriceHandler {
	return &GetProductWithPriceHandler{
		queryService: queryService,
	}
}

func (h *GetProductWithPriceHandler) Handle(ctx context.Context, id int64) (*ProductWithPriceDTO, error) {
	return h.queryService.GetWithCurrentPrice(ctx, id)
}
