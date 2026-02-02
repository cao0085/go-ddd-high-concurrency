package product

import (
	"context"
	"database/sql"

	domain "flash-sale-order-system/internal/domain/product"
)

type RemoveProductCommand struct {
	Id int64
}

type RemoveProductHandler struct {
	db          *sql.DB
	productRepo domain.ProductRepository
}

func NewRemoveProductHandler(
	db *sql.DB,
	productRepo domain.ProductRepository,
) *RemoveProductHandler {
	return &RemoveProductHandler{
		db:          db,
		productRepo: productRepo,
	}
}

func (h *RemoveProductHandler) Handle(ctx context.Context, cmd RemoveProductCommand) error {

	product, err := h.productRepo.FindByID(ctx, cmd.Id)
	if err != nil {
		return err
	}

	if err := product.CanDelete(); err != nil {
		return err
	}

	if err := h.productRepo.Delete(ctx, product.ID()); err != nil {
		return err
	}

	return nil
}
