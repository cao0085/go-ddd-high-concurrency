package product

import (
	"errors"
	"fmt"

	shareddomain "flash-sale-order-system/internal/shared/domain"
)

// PriceList is a value object that holds multiple currency prices for a product
type PriceList struct {
	prices map[shareddomain.Currency]shareddomain.Money // key: currency, value: money
}

// NewPriceList creates a new PriceList with at least one price
func NewPriceList(prices map[shareddomain.Currency]shareddomain.Money) (PriceList, error) {
	if len(prices) == 0 {
		return PriceList{}, errors.New("price list must have at least one price")
	}

	// Create a copy to ensure immutability
	pricesCopy := make(map[shareddomain.Currency]shareddomain.Money, len(prices))
	for currency, money := range prices {
		pricesCopy[currency] = money
	}

	return PriceList{prices: pricesCopy}, nil
}

// GetPrice returns the price for a specific currency
func (pl PriceList) GetPrice(currency shareddomain.Currency) (shareddomain.Money, error) {
	price, exists := pl.prices[currency]
	if !exists {
		return shareddomain.Money{}, fmt.Errorf("price not available for currency: %s", currency)
	}
	return price, nil
}

// GetAllPrices returns a copy of all prices
func (pl PriceList) GetAllPrices() map[shareddomain.Currency]shareddomain.Money {
	pricesCopy := make(map[shareddomain.Currency]shareddomain.Money, len(pl.prices))
	for currency, money := range pl.prices {
		pricesCopy[currency] = money
	}
	return pricesCopy
}