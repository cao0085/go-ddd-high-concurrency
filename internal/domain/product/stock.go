package product

import (
	"errors"
	"math"
)

// Value Object
type Stock struct {
	available int32
	reserved  int32
}

func Initial(available int32) (Stock, error) {
	if available < 0 {
		return Stock{}, errors.New("stock cannot be negative")
	}
	return Stock{available: available, reserved: 0}, nil
}

func (s Stock) Add(quantity int32) (Stock, error) {
	if quantity < 0 {
		return s, errors.New("quantity to add cannot be negative")
	}
	if quantity > 100000 || s.available > math.MaxInt32-quantity {
		return s, errors.New("stock overflow")
	}
	return Stock{
		available: s.available + quantity,
		reserved:  s.reserved,
	}, nil
}

func (s Stock) Reserve(quantity int32) (Stock, error) {
	if quantity <= 0 {
		return s, errors.New("quantity must be positive")
	}
	if s.available < quantity {
		return s, errors.New("insufficient stock")
	}
	return Stock{
		available: s.available - quantity,
		reserved:  s.reserved + quantity,
	}, nil
}

func (s Stock) ConfirmReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity {
		return s, errors.New("insufficient reserved stock to confirm")
	}

	return Stock{
		available: s.available,
		reserved:  s.reserved - quantity,
	}, nil
}

func (s Stock) CancelReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity {
		return s, errors.New("insufficient reserved stock to cancel")
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
		return Stock{}, errors.New("insufficient stock for adjustment")
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
