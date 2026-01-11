// internal/application/product/reserve_stock.go
package product

import (
	"context"
	"fmt"

	domain "flash-sale-order-system/internal/domain/product"
)

type ReserveStockCommand struct {
	ProductID int64
	Quantity  int32
}

type ReserveStockHandler struct {
	productRepo domain.Repository
}

func NewReserveStockHandler(repo domain.Repository) *ReserveStockHandler {
	return &ReserveStockHandler{productRepo: repo}
}

func (h *ReserveStockHandler) Handle(ctx context.Context, cmd ReserveStockCommand) error {
	// 1. 取得產品
	product, err := h.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	// 2. 呼叫業務邏輯：預留庫存
	newStock, err := product.Stock().Reserve(cmd.Quantity)
	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// 3. 持久化新的庫存狀態
	if err := h.productRepo.UpdateStock(ctx, cmd.ProductID, newStock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
