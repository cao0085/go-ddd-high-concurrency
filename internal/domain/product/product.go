package product

import "time"

// Aggregate
type Product struct {
	id          int64
	sku         string
	name        string
	description string
	status      int8
	createdAt   time.Time
	updatedAt   time.Time
	stock       Stock
}

// NewProduct creates a new product with a given ID
func NewProduct(id int64, name string, description string, sku string, quantity int32) (*Product, error) {
	if name == "" {
		return nil, ErrEmptyProductName
	}
	if sku == "" {
		return nil, ErrEmptySKU
	}

	stock, err := NewStock(quantity)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Product{
		id:          id,
		name:        name,
		description: description,
		sku:         sku,
		status:      StatusInactive,
		createdAt:   now,
		updatedAt:   now,
		stock:       stock,
	}, nil
}

func (p *Product) UpdateInfo(name string, description string, status int8) error {
	if name == "" {
		return ErrEmptyProductName
	}

	p.name = name
	p.status = status
	p.description = description
	p.updatedAt = time.Now()

	return nil
}

func (p *Product) Activate() error {
	if p.status == StatusActive {
		return ErrAlreadyActive
	}
	p.status = StatusActive
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) Deactivate() error {
	if p.status == StatusInactive {
		return ErrAlreadyInactive
	}
	p.status = StatusInactive
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) IsActive() bool {
	return p.status == StatusActive
}

func (p *Product) CanDelete() error {
	if p.stock.Reserved() > 0 {
		return ErrHasReservedStock
	}
	// onGoing: add more business rules if needed
	return nil
}

// Getters
func (p *Product) ID() int64            { return p.id }
func (p *Product) SKU() string          { return p.sku }
func (p *Product) Name() string         { return p.name }
func (p *Product) Description() string  { return p.description }
func (p *Product) Status() int8         { return p.status }
func (p *Product) Stock() Stock         { return p.stock }
func (p *Product) CreatedAt() time.Time { return p.createdAt }
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }
