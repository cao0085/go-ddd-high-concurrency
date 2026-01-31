package product

import "math"

// Value Object
type Stock struct {
	available int32
	reserved  int32
}

func NewStock(available int32) (Stock, error) {
	if available < 0 {
		return Stock{}, ErrNegativeStock
	}
	return Stock{available: available, reserved: 0}, nil
}

func (s Stock) Add(quantity int32) (Stock, error) {
	if quantity < 0 {
		return s, ErrNegativeQuantity
	}
	if quantity > 100000 || s.available > math.MaxInt32-quantity {
		return s, ErrStockOverflow
	}
	return Stock{
		available: s.available + quantity,
		reserved:  s.reserved,
	}, nil
}

func (s Stock) Reserve(quantity int32) (Stock, error) {
	if quantity <= 0 {
		return s, ErrNonPositiveQuantity
	}
	if s.available < quantity {
		return s, ErrInsufficientStock
	}
	return Stock{
		available: s.available - quantity,
		reserved:  s.reserved + quantity,
	}, nil
}

func (s Stock) ConfirmReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity {
		return s, ErrInsufficientReserved
	}

	return Stock{
		available: s.available,
		reserved:  s.reserved - quantity,
	}, nil
}

func (s Stock) CancelReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity {
		return s, ErrInsufficientReserved
	}

	return Stock{
		available: s.available + quantity,
		reserved:  s.reserved - quantity,
	}, nil
}

// positive or negative
func (s Stock) AdjustAvailable(delta int32) (Stock, error) {
	newAvailable := s.available + delta

	if newAvailable < 0 {
		return Stock{}, ErrInsufficientStock
	}

	return Stock{
		available: newAvailable,
		reserved:  s.reserved,
	}, nil
}

// getter methods
func (s Stock) Available() int32 { return s.available }
func (s Stock) Reserved() int32  { return s.reserved }
func (s Stock) Total() int32     { return s.available + s.reserved }
