package command

import (
	"context"
	"database/sql"

	domain "flash-sale-order-system/internal/domain/product"
)

type UpdateProductInfoCommand struct {
	Id          int64
	Name        string
	Description string
	Status      int8
}

type UpdateProductInfoHandler struct {
	db          *sql.DB
	productRepo domain.ProductRepository
}

func NewUpdateProductInfoHandler(
	db *sql.DB,
	productRepo domain.ProductRepository,
) *UpdateProductInfoHandler {
	return &UpdateProductInfoHandler{
		db:          db,
		productRepo: productRepo,
	}
}

func (h *UpdateProductInfoHandler) Handle(ctx context.Context, cmd UpdateProductInfoCommand) error {

	product, err := h.productRepo.FindByID(ctx, cmd.Id)
	if err != nil {
		return err
	}

	if err := product.UpdateInfo(cmd.Name, cmd.Description, cmd.Status); err != nil {
		return err
	}

	if err := h.productRepo.UpdateInfo(ctx, product); err != nil {
		return err
	}

	return nil

}
