package product

import (
	"errors"
	"time"
)

// Product represents a product entity in the system.
type Product struct {
	id        int64
	name      string
	stock     stock
	priceList PriceList
	createdAt time.Time
	updatedAt time.Time
}

// NewProduct creates a new product
func NewProduct(name string, stock stock, priceList PriceList) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	if priceList.IsEmpty() {
		return nil, errors.New("product must have at least one price")
	}

	now := time.Now()
	return &Product{
		name:      name,
		stock:     stock,
		priceList: priceList,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// UpdateInfo updates the product's name and price list
func (p *Product) UpdateInfo(name string, priceList PriceList) error {
    if name == "" {
        return errors.New("product name cannot be empty")
    }
    
    if priceList.IsEmpty() {
        return errors.New("price list cannot be empty")
    }
    
    p.name = name
    p.priceList = priceList
    p.updatedAt = time.Now()
    
    return nil
}

// CanDelete checks if the product can be deleted based on business rules.
func (p *Product) CanDelete() error {

    if p.stock.Reserved() > 0 {
        return errors.New("cannot delete product with reserved stock")
    }

	// onGoing: add more business rules if needed
    
    return nil
}

// Getters
func (p *Product) ID() int64                       { return p.id }
func (p *Product) Name() string                    { return p.name }
func (p *Product) Stock() int32                    { return p.stock }
func (p *Product) PriceList() PriceList            { return p.priceList }
func (p *Product) CreatedAt() time.Time            { return p.createdAt }
func (p *Product) UpdatedAt() time.Time            { return p.updatedAt }

// GetPrice returns the price for a specific currency
func (p *Product) GetPrice(currency shareddomain.Currency) (shareddomain.Money, error) {
	return p.priceList.GetPrice(currency)
}