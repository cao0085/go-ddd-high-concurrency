// scheduled_price.go
package product

import (
    "time"
    shareddomain "flash-sale-order-system/internal/shared/domain"
)

// ScheduledPrice entity - 價格排程的實體
type ScheduledPrice struct {
    id         string
    prices     map[shareddomain.Currency]shareddomain.Money  // 不用 PriceList
    validFrom  time.Time
    validUntil *time.Time  // nil = 永久有效
    createdAt  time.Time
}

func newScheduledPrice(
    prices map[shareddomain.Currency]shareddomain.Money,
    from time.Time,
    until *time.Time,
) *ScheduledPrice {
    // 複製 prices 確保不可變性
    pricesCopy := make(map[shareddomain.Currency]shareddomain.Money)
    for k, v := range prices {
        pricesCopy[k] = v
    }
    
    return &ScheduledPrice{
        id:         generateID(),  // 你的 ID 生成邏輯
        prices:     pricesCopy,
        validFrom:  from,
        validUntil: until,
        createdAt:  time.Now(),
    }
}

func (sp *ScheduledPrice) isValidAt(t time.Time) bool {
    if t.Before(sp.validFrom) {
        return false
    }
    
    if sp.validUntil != nil && t.After(*sp.validUntil) {
        return false
    }
    
    return true
}

func (sp *ScheduledPrice) overlaps(from time.Time, until *time.Time) bool {
    // 時間區間重疊邏輯
    // [from1, until1] overlaps [from2, until2]
    
    // 如果 until 是 nil，視為無限大
    end1 := sp.validUntil
    end2 := until
    
    // 實作重疊檢查邏輯...
    // 簡化版：
    if from.After(sp.validFrom) || from.Equal(sp.validFrom) {
        if end1 == nil || from.Before(*end1) {
            return true
        }
    }
    // ... 完整邏輯
    
    return false
}

func (sp *ScheduledPrice) getCopyOfPrices() map[shareddomain.Currency]shareddomain.Money {
    copy := make(map[shareddomain.Currency]shareddomain.Money)
    for k, v := range sp.prices {
        copy[k] = v
    }
    return copy
}

// Getters
func (sp *ScheduledPrice) ID() string { return sp.id }
func (sp *ScheduledPrice) ValidFrom() time.Time { return sp.validFrom }
func (sp *ScheduledPrice) ValidUntil() *time.Time { return sp.validUntil }
