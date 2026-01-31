package product

import (
	"errors"
	"fmt"

	shareddomain "flash-sale-order-system/internal/shared/domain"
)

// Value Object
// Example: A product costs USD 10.00, TWD 300, JPY 1500
type MultiCurrencyPrice struct {
	prices map[shareddomain.Currency]shareddomain.Money
}

// Create a copy to ensure immutability
func NewMultiCurrencyPrice(prices map[shareddomain.Currency]shareddomain.Money) (MultiCurrencyPrice, error) {
	if len(prices) == 0 {
		return MultiCurrencyPrice{}, errors.New("multi currency price must have at least one price")
	}

	pricesCopy := make(map[shareddomain.Currency]shareddomain.Money, len(prices))
	for currency, money := range prices {
		pricesCopy[currency] = money
	}

	return MultiCurrencyPrice{prices: pricesCopy}, nil
}

func (p MultiCurrencyPrice) GetPrice(currency shareddomain.Currency) (shareddomain.Money, error) {
	price, exists := p.prices[currency]
	if !exists {
		return shareddomain.Money{}, fmt.Errorf("price not available for currency: %s", currency)
	}
	return price, nil
}

func (p MultiCurrencyPrice) GetAllPrices() map[shareddomain.Currency]shareddomain.Money {
	pricesCopy := make(map[shareddomain.Currency]shareddomain.Money, len(p.prices))
	for currency, money := range p.prices {
		pricesCopy[currency] = money
	}
	return pricesCopy
}
