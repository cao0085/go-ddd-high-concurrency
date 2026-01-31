package product

import "errors"

// Stock errors
var (
	ErrNegativeStock        = errors.New("stock cannot be negative")
	ErrNegativeQuantity     = errors.New("quantity cannot be negative")
	ErrNonPositiveQuantity  = errors.New("quantity must be positive")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInsufficientReserved = errors.New("insufficient reserved stock")
	ErrStockOverflow        = errors.New("stock overflow")
)

// Product errors
var (
	ErrEmptyProductName      = errors.New("product name cannot be empty")
	ErrEmptySKU              = errors.New("product SKU cannot be empty")
	ErrProductNotFound       = errors.New("product not found")
	ErrHasReservedStock      = errors.New("cannot delete product with reserved stock")
	ErrAlreadyActive         = errors.New("product is already active")
	ErrAlreadyInactive       = errors.New("product is already inactive")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// Status constants
const (
	StatusActive   int8 = 1
	StatusInactive int8 = 9
)

// Pricing errors
var (
	ErrPeriodOverlap       = errors.New("price period overlaps with existing period")
	ErrInvalidPeriod       = errors.New("invalid period: end time is before start time")
	ErrNoPriceFound        = errors.New("no valid price found for the given time")
	ErrCurrencyNotFound    = errors.New("price not available for currency")
	ErrEmptyMultiCurrency  = errors.New("multi currency price must have at least one price")
)
