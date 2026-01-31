package product

import (
	"context"
	"database/sql"
	"time"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	"flash-sale-order-system/internal/Infrastructure/persistence/tx"
	domain "flash-sale-order-system/internal/domain/product"
	shareddomain "flash-sale-order-system/internal/shared/domain"
)

type CreateProductCommand struct {
	Name        string
	Description string
	SKU         string
	Quantity    int32
	Prices      map[shareddomain.Currency]shareddomain.Money
	PriceFrom   time.Time
	PriceUntil  *time.Time
}

type CreateProductHandler struct {
	db          *sql.DB
	idGenerator *idgen.IDGenerator
	productRepo domain.ProductRepository
	pricingRepo domain.ProductPricingRepository
}

func NewCreateProductHandler(
	db *sql.DB,
	idGen *idgen.IDGenerator,
	productRepo domain.ProductRepository,
	pricingRepo domain.ProductPricingRepository,
) *CreateProductHandler {
	return &CreateProductHandler{
		db:          db,
		idGenerator: idGen,
		productRepo: productRepo,
		pricingRepo: pricingRepo,
	}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (int64, error) {
	// 1. snowflake id generator
	productID := h.idGenerator.Generate()

	// 2. Product Aggregate
	product, err := domain.NewProduct(
		productID,
		cmd.Name,
		cmd.Description,
		cmd.SKU,
		cmd.Quantity,
	)
	if err != nil {
		return 0, err
	}

	// 3. ProductPricing Aggregate
	pricing := domain.NewProductPricing(productID)

	prices, err := domain.NewMultiCurrencyPrice(cmd.Prices)
	if err != nil {
		return 0, err
	}

	if err := pricing.AddPeriod(prices, cmd.PriceFrom, cmd.PriceUntil); err != nil {
		return 0, err
	}

	// 4. Transactional Save
	err = tx.WithTx(ctx, h.db, func(txCtx context.Context) error {
		if err := h.productRepo.Insert(txCtx, product); err != nil {
			return err
		}
		if err := h.pricingRepo.Save(txCtx, pricing); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return productID, nil
}
