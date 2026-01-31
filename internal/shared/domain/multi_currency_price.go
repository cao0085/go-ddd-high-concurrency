package domain

// Value Object
// Example: A product costs USD 10.00, TWD 300, JPY 1500
type MultiCurrencyPrice struct {
	prices map[Currency]Money
}

// Create a copy to ensure immutability
func NewMultiCurrencyPrice(prices map[Currency]Money) (MultiCurrencyPrice, error) {
	if len(prices) == 0 {
		return MultiCurrencyPrice{}, ErrEmptyMultiCurrency
	}

	pricesCopy := make(map[Currency]Money, len(prices))
	for currency, money := range prices {
		pricesCopy[currency] = money
	}

	return MultiCurrencyPrice{prices: pricesCopy}, nil
}

func (p MultiCurrencyPrice) GetPrice(currency Currency) (Money, error) {
	price, exists := p.prices[currency]
	if !exists {
		return Money{}, ErrCurrencyNotFound
	}
	return price, nil
}

func (p MultiCurrencyPrice) GetAllPrices() map[Currency]Money {
	pricesCopy := make(map[Currency]Money, len(p.prices))
	for currency, money := range p.prices {
		pricesCopy[currency] = money
	}
	return pricesCopy
}

func (p MultiCurrencyPrice) Currencies() []Currency {
	currencies := make([]Currency, 0, len(p.prices))
	for c := range p.prices {
		currencies = append(currencies, c)
	}
	return currencies
}
