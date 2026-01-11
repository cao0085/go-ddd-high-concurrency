// internal/application/product/cancel_reservation.go
package product

import (
	"context"
	"fmt"

	domain "flash-sale-order-system/internal/domain/product"
)

type CancelReservationCommand struct {
	ProductID int64
	Quantity  int32
}

type CancelReservationHandler struct {
	productRepo domain.Repository
}

func NewCancelReservationHandler(repo domain.Repository) *CancelReservationHandler {
	return &CancelReservationHandler{productRepo: repo}
}

func (h *CancelReservationHandler) Handle(ctx context.Context, cmd CancelReservationCommand) error {
	// 1. 取得產品
	product, err := h.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	// 2. 呼叫業務邏輯：取消預留
	newStock, err := product.Stock().CancelReservation(cmd.Quantity)
	if err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	// 3. 持久化新的庫存狀態
	if err := h.productRepo.UpdateStock(ctx, cmd.ProductID, newStock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
