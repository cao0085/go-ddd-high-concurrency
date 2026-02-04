package product

import "context"

type ProductRepository interface {
	Insert(ctx context.Context, p *Product) error
	UpdateInfo(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Product, error)
}

type ProductPricingRepository interface {
	FindByProductID(ctx context.Context, productID int64) (*ProductPricing, error)
	Save(ctx context.Context, productPricing *ProductPricing) error
}
