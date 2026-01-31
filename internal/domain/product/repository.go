package product

import "context"

// Repository is the interface for Product aggregate
type Repository interface {
	// 基本 CRUD
	FindByID(ctx context.Context, id int64) (*Product, error)
	Save(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error

	// 查詢方法
	FindByIDs(ctx context.Context, ids []int64) ([]*Product, error)
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// 庫存操作
	UpdateStock(ctx context.Context, productID int64, newStock Stock) error
}

// PricingRepository is the interface for ProductPricing aggregate
type PricingRepository interface {
	FindByProductID(ctx context.Context, productID int64) (*ProductPricing, error)
	Save(ctx context.Context, pricing *ProductPricing) error
	Delete(ctx context.Context, productID int64) error
}