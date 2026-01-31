package domain

import "errors"

var (
	ErrEmptyMultiCurrency = errors.New("multi currency price cannot be empty")
	ErrCurrencyNotFound   = errors.New("currency not found")
)
