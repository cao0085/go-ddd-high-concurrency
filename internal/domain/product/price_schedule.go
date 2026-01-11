// price_schedule.go
package product

import (
    "errors"
    "time"
    shareddomain "flash-sale-order-system/internal/shared/domain"
)

// PriceSchedule aggregate root - 管理價格排程
type PriceSchedule struct {
    productID int64
    schedules []*ScheduledPrice
}

func NewPriceSchedule(productID int64) *PriceSchedule {
    return &PriceSchedule{
        productID: productID,
        schedules: []*ScheduledPrice{},
    }
}

func (ps *PriceSchedule) AddSchedule(
    prices map[shareddomain.Currency]shareddomain.Money,
    from time.Time,
    until *time.Time,
) error {
    // 檢查重疊
    if ps.hasOverlap(from, until) {
        return errors.New("schedule overlaps with existing period")
    }
    
    schedule := newScheduledPrice(prices, from, until)
    ps.schedules = append(ps.schedules, schedule)
    
    return nil
}

func (ps *PriceSchedule) hasOverlap(from time.Time, until *time.Time) bool {
    for _, s := range ps.schedules {
        if s.overlaps(from, until) {
            return true
        }
    }
    return false
}

func (ps *PriceSchedule) GetCurrentPrices(now time.Time) (map[shareddomain.Currency]shareddomain.Money, error) {
    for _, s := range ps.schedules {
        if s.isValidAt(now) {
            return s.getCopyOfPrices(), nil
        }
    }
    return nil, errors.New("no valid price found")
}

// Getters
func (ps *PriceSchedule) ProductID() int64 { return ps.productID }
func (ps *PriceSchedule) Schedules() []*ScheduledPrice { return ps.schedules }
