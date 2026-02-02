package product

import (
	"context"
	"database/sql"
	"time"

	domain "flash-sale-order-system/internal/domain/product"
	shareddomain "flash-sale-order-system/internal/shared/domain"
)

type SaveProductPricesCommand struct {
	ProductID int64
	Periods   []PricePeriodInput
}

type PricePeriodInput struct {
	Currency   string
	Amount     float64
	ValidFrom  time.Time
	ValidUntil *time.Time
}

type SaveProductPricesHandler struct {
	db         *sql.DB
	pricesRepo domain.ProductPricingRepository
}

func NewSaveProductPricesHandler(
	db *sql.DB,
	pricesRepo domain.ProductPricingRepository,
) *SaveProductPricesHandler {
	return &SaveProductPricesHandler{
		db:         db,
		pricesRepo: pricesRepo,
	}
}

func (h *SaveProductPricesHandler) Handle(ctx context.Context, cmd SaveProductPricesCommand) error {

	pp, err := h.pricesRepo.FindByProductID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	for _, p := range cmd.Periods {
		price, err := shareddomain.NewSinglePrice(p.Amount, shareddomain.Currency(p.Currency))
		if err != nil {
			return err
		}
		if err := pp.AddPeriod(price, p.ValidFrom, p.ValidUntil); err != nil {
			return err
		}
	}

	if err := h.pricesRepo.Save(ctx, pp); err != nil {
		return err
	}

	return nil
}
