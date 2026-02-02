package product

import (
	"context"
	"database/sql"

	domain "flash-sale-order-system/internal/domain/product"
)

type GetProductCommand struct {
	Id int64
}

type GetProductHandler struct {
	db          *sql.DB
	productRepo domain.ProductRepository
}

func NewGetProductHandler(
	db *sql.DB,
	productRepo domain.ProductRepository,
) *GetProductHandler {
	return &GetProductHandler{
		db:          db,
		productRepo: productRepo,
	}
}

func (h *GetProductHandler) Handle(ctx context.Context, id int64) (*domain.Product, error) {

	product, err := h.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
