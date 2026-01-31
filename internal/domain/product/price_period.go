// price_period.go
package product

import "time"

// Value Object
type PricePeriod struct {
	prices     MultiCurrencyPrice
	validFrom  time.Time
	validUntil *time.Time
}

func NewPricePeriod(
	prices MultiCurrencyPrice,
	from time.Time,
	until *time.Time,
) (PricePeriod, error) {

	if until != nil && until.Before(from) {
		return PricePeriod{}, ErrInvalidPeriod
	}

	return PricePeriod{
		prices:     prices,
		validFrom:  from,
		validUntil: until,
	}, nil
}

func (p PricePeriod) IsValidAt(t time.Time) bool {
	if t.Before(p.validFrom) {
		return false
	}

	if p.validUntil != nil && t.After(*p.validUntil) {
		return false
	}

	return true
}

func (p PricePeriod) Overlaps(from time.Time, until *time.Time, prices MultiCurrencyPrice) bool {
	// 1. 先檢查是否有相同幣別
	if !p.hasCommonCurrency(prices) {
		return false // 沒有相同幣別，不算重疊
	}

	// 2. 再檢查時間重疊
	// 區間 A: [p.validFrom, p.validUntil]
	// 區間 B: [from, until]
	// 重疊條件: A.start <= B.end AND B.start <= A.end

	// 檢查 B.start <= A.end
	if p.validUntil != nil && from.After(*p.validUntil) {
		return false
	}

	// 檢查 A.start <= B.end
	if until != nil && p.validFrom.After(*until) {
		return false
	}

	return true
}

func (p PricePeriod) hasCommonCurrency(other MultiCurrencyPrice) bool {
	for currency := range other.GetAllPrices() {
		if _, err := p.prices.GetPrice(currency); err == nil {
			return true // 找到相同幣別
		}
	}
	return false
}

// Getters
func (p PricePeriod) Prices() MultiCurrencyPrice { return p.prices }
func (p PricePeriod) ValidFrom() time.Time       { return p.validFrom }
func (p PricePeriod) ValidUntil() *time.Time     { return p.validUntil }
