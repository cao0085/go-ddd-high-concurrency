package query

import "context"

// ProductQueryService defines the interface for product read operations.
// This interface is implemented in the Infrastructure layer.
type ProductQueryService interface {
	GetByID(ctx context.Context, id int64) (*ProductDTO, error)
	GetWithCurrentPrice(ctx context.Context, id int64) (*ProductWithPriceDTO, error)
}
