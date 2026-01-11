// internal/application/product/create_product.go
package product

import (
    "context"
    
    domain "flash-sale-order-system/internal/domain/product"
    shared "flash-sale-order-system/internal/shared/domain"
)

type DeleteProductCommand struct {
    ID           string
}

type DeleteProductHandler struct {
    productRepo domain.Repository
}

func NewDeleteProductHandler(repo domain.Repository) *DeleteProductHandler {
    return &DeleteProductHandler{productRepo: repo}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
    return h.productRepo.Delete(ctx, cmd.ID)
}
