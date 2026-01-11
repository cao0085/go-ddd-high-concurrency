// internal/application/product/confirm_reservation.go
package product

import (
	"context"
	"fmt"

	domain "flash-sale-order-system/internal/domain/product"
)

type ConfirmReservationCommand struct {
	ProductID int64
	Quantity  int32
}

type ConfirmReservationHandler struct {
	productRepo domain.Repository
}

func NewConfirmReservationHandler(repo domain.Repository) *ConfirmReservationHandler {
	return &ConfirmReservationHandler{productRepo: repo}
}

func (h *ConfirmReservationHandler) Handle(ctx context.Context, cmd ConfirmReservationCommand) error {
	// 1. 取得產品
	product, err := h.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	// 2. 呼叫業務邏輯：確認預留
	newStock, err := product.Stock().ConfirmReservation(cmd.Quantity)
	if err != nil {
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	// 3. 持久化新的庫存狀態
	if err := h.productRepo.UpdateStock(ctx, cmd.ProductID, newStock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
