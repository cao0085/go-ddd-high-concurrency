// internal/application/product/create_product.go
package product

import (
	"context"
	"time"

	domain "flash-sale-order-system/internal/domain/product"
	"flash-sale-order-system/internal/Infrastructure/idgen"
	shareddomain "flash-sale-order-system/internal/shared/domain"
)

type CreateProductCommand struct {
	Name        string
	Description string
	SKU         string
	Quantity    int32
	Prices      map[shareddomain.Currency]shareddomain.Money // 價格
	PriceFrom   time.Time                                    // 價格生效時間
	PriceUntil  *time.Time                                   // 價格結束時間（nil = 永久）
}

type CreateProductHandler struct {
	idGenerator *idgen.IDGenerator
	productRepo domain.Repository
	pricingRepo domain.PricingRepository
}

func NewCreateProductHandler(
	idGen *idgen.IDGenerator,
	productRepo domain.Repository,
	pricingRepo domain.PricingRepository,
) *CreateProductHandler {
	return &CreateProductHandler{
		idGenerator: idGen,
		productRepo: productRepo,
		pricingRepo: pricingRepo,
	}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (int64, error) {
	// 1. 產生 ID
	productID := h.idGenerator.Generate()

	// 2. 建立 Product
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

	// 3. 建立 ProductPricing
	pricing := domain.NewProductPricing(productID)

	prices, err := domain.NewMultiCurrencyPrice(cmd.Prices)
	if err != nil {
		return 0, err
	}

	if err := pricing.AddPeriod(prices, cmd.PriceFrom, cmd.PriceUntil); err != nil {
		return 0, err
	}

	// 4. 儲存 Product
	if err := h.productRepo.Save(ctx, product); err != nil {
		return 0, err
	}

	// 5. 儲存 Pricing
	if err := h.pricingRepo.Save(ctx, pricing); err != nil {
		return 0, err
	}

	return productID, nil
}
