// product_pricing.go
package product

import (
	"time"

	shareddomain "flash-sale-order-system/internal/shared/domain"
)

// Aggregate
type ProductPricing struct {
	productID int64
	periods   []PricePeriod
}

func NewProductPricing(productID int64) *ProductPricing {
	return &ProductPricing{
		productID: productID,
		periods:   []PricePeriod{},
	}
}

func (pp *ProductPricing) AddPeriod(
	prices shareddomain.MultiCurrencyPrice,
	from time.Time,
	until *time.Time,
) error {

	period, err := NewPricePeriod(prices, from, until)
	if err != nil {
		return err
	}

	if pp.hasOverlap(from, until, prices) {
		return ErrPeriodOverlap
	}

	pp.periods = append(pp.periods, period)
	return nil
}

func (pp *ProductPricing) GetCurrentPrices(now time.Time) (shareddomain.MultiCurrencyPrice, error) {
	for _, period := range pp.periods {
		if period.IsValidAt(now) {
			return period.Prices(), nil
		}
	}
	return shareddomain.MultiCurrencyPrice{}, ErrNoPriceFound
}

func (pp *ProductPricing) GetPriceForCurrency(now time.Time, currency shareddomain.Currency) (shareddomain.Money, error) {
	prices, err := pp.GetCurrentPrices(now)
	if err != nil {
		return shareddomain.Money{}, err
	}
	return prices.GetPrice(currency)
}

func (pp *ProductPricing) hasOverlap(from time.Time, until *time.Time, prices shareddomain.MultiCurrencyPrice) bool {
	for _, period := range pp.periods {
		if period.Overlaps(from, until, prices) {
			return true
		}
	}
	return false
}

// Getters
func (pp *ProductPricing) ProductID() int64       { return pp.productID }
func (pp *ProductPricing) Periods() []PricePeriod { return pp.periods }
