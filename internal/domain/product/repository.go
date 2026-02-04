package product

import "context"

type ProductRepository interface {
	Insert(ctx context.Context, p *Product) error
	UpdateInfo(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductPricingRepository interface {
	Save(ctx context.Context, productPricing *ProductPricing) error
}
