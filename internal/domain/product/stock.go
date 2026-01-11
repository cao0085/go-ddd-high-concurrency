package product

import (
    "errors"
    "math"
)

type Stock struct {
    available int32  // 可用庫存
    reserved  int32  // 已預留(鎖定)庫存
}

// NewStock creates a new Stock instance with the given available quantity.
func NewStock(available int32) (Stock, error) {
    if available < 0 {
        return Stock{}, errors.New("stock cannot be negative")
    }
    return Stock{available: available, reserved: 0}, nil
}

// AddStock adds the specified quantity to available stock.
func (s *Stock) AddStock(quantity int32) (Stock, error) {
    if quantity < 0 {
        return s, errors.New("quantity to add cannot be negative")
    }
    if quantity > math.MaxInt32 - s.available {
        return s, errors.New("stock overflow")
    }
    return Stock{
        available: s.available + quantity,
        reserved:  s.reserved,
    }, nil
}

// Reserve reserves the specified quantity from available stock.
func (s *Stock) Reserve(quantity int32) (Stock, error) {
    if quantity <= 0 { // Input validation
        return s, errors.New("quantity must be positive")
    }
    if s.available < quantity { // State validation
        return s, errors.New("insufficient stock")
    }
    return Stock{
        available: s.available - quantity,
        reserved:  s.reserved + quantity,
    }, nil
}

// ConfirmReservation confirms the reservation, reducing the reserved stock.
func (s *Stock) ConfirmReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity {
		return s, errors.New("insufficient reserved stock to confirm")
	}

	return Stock{
		available: s.available,
		reserved:  s.reserved - quantity,
	}, nil
}

// CancelReservation cancels the reservation, returning the reserved stock to available.
func (s *Stock) CancelReservation(quantity int32) (Stock, error) {
	if s.reserved < quantity { // 防止重複取消
		return s, errors.New("insufficient reserved stock to cancel")
	}

	return Stock{
		available: s.available + quantity,
		reserved:  s.reserved - quantity,
	}, nil
}

// AdjustAvailable adjusts the available stock by the given delta.
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
