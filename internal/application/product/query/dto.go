package query

import "time"

type ProductDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Status      int8      `json:"status"`
	Stock       StockDTO  `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StockDTO struct {
	Available int32 `json:"available"`
	Reserved  int32 `json:"reserved"`
}

type ProductWithPriceDTO struct {
	ProductDTO
	CurrentPrice *PriceDTO `json:"current_price,omitempty"`
}

type PriceDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
