package domain

import (
	"errors"
	"fmt"
	"math"
)

type Currency string

const (
	USD Currency = "USD"
	TWD Currency = "TWD"
	JPY Currency = "JPY"
)

type Money struct {
	amount   float64
	currency Currency
}

func NewMoney(amount float64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	// 根據幣別驗證精度
	if err := validatePrecision(amount, currency); err != nil {
		return Money{}, err
	}

	return Money{amount: amount, currency: currency}, nil
}

func validatePrecision(amount float64, currency Currency) error {
	precision := getPrecision(currency)
	rounded := round(amount, precision)

	if amount != rounded {
		return fmt.Errorf("%s can only have %d decimal places", currency, precision)
	}
	return nil
}

func getPrecision(c Currency) int {
	switch c {
	case TWD, JPY:
		return 0 // 台幣/日圓: 1000
	case USD:
		return 2 // 美金: 100.50
	default:
		return 2
	}
}

func (m Money) Amount() float64    { return m.amount }
func (m Money) Currency() Currency { return m.currency }

func (m Money) ConvertTo(target Currency, exchangeRate float64) Money {
	newAmount := m.amount * exchangeRate
	precision := getPrecision(target)

	return Money{
		amount:   round(newAmount, precision),
		currency: target,
	}
}

func round(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}
