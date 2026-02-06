package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// StockCache handles stock-related Redis operations
type StockCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewStockCache creates a new StockCache instance
func NewStockCache(client *redis.Client) *StockCache {
	return &StockCache{
		client: client,
		ttl:    24 * time.Hour, // 預設 TTL 24 小時
	}
}

// stockKey generates Redis key for product stock
func (s *StockCache) stockKey(productID int64) string {
	return fmt.Sprintf("stock:product:%d", productID)
}

// availableKey generates Redis key for available stock
func (s *StockCache) availableKey(productID int64) string {
	return fmt.Sprintf("stock:product:%d:available", productID)
}

// reservedKey generates Redis key for reserved stock
func (s *StockCache) reservedKey(productID int64) string {
	return fmt.Sprintf("stock:product:%d:reserved", productID)
}

// InitStock initializes stock in Redis (從資料庫同步)
func (s *StockCache) InitStock(ctx context.Context, productID int64, available, reserved int32) error {
	pipe := s.client.Pipeline()

	availKey := s.availableKey(productID)
	reservKey := s.reservedKey(productID)

	pipe.Set(ctx, availKey, available, s.ttl)
	pipe.Set(ctx, reservKey, reserved, s.ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to init stock: %w", err)
	}

	return nil
}

// GetAvailable gets available stock from Redis
func (s *StockCache) GetAvailable(ctx context.Context, productID int64) (int32, error) {
	val, err := s.client.Get(ctx, s.availableKey(productID)).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("stock not found in cache")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get available stock: %w", err)
	}

	stock, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid stock value: %w", err)
	}

	return int32(stock), nil
}

// Reserve reserves stock atomically using Lua script
// 返回 true 表示預扣成功，false 表示庫存不足
func (s *StockCache) Reserve(ctx context.Context, productID int64, quantity int32) (bool, error) {
	// Lua script 保證原子性
	script := `
		local availKey = KEYS[1]
		local reservKey = KEYS[2]
		local quantity = tonumber(ARGV[1])
		
		local available = tonumber(redis.call('GET', availKey) or 0)
		
		if available < quantity then
			return 0
		end
		
		redis.call('DECRBY', availKey, quantity)
		redis.call('INCRBY', reservKey, quantity)
		return 1
	`

	result, err := s.client.Eval(ctx, script, []string{
		s.availableKey(productID),
		s.reservedKey(productID),
	}, quantity).Int()

	if err != nil {
		return false, fmt.Errorf("failed to reserve stock: %w", err)
	}

	return result == 1, nil
}

// ConfirmReservation confirms a reservation (扣除 reserved)
func (s *StockCache) ConfirmReservation(ctx context.Context, productID int64, quantity int32) error {
	script := `
		local reservKey = KEYS[1]
		local quantity = tonumber(ARGV[1])
		
		local reserved = tonumber(redis.call('GET', reservKey) or 0)
		
		if reserved < quantity then
			return 0
		end
		
		redis.call('DECRBY', reservKey, quantity)
		return 1
	`

	result, err := s.client.Eval(ctx, script, []string{
		s.reservedKey(productID),
	}, quantity).Int()

	if err != nil {
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	if result == 0 {
		return fmt.Errorf("insufficient reserved stock")
	}

	return nil
}

// CancelReservation cancels a reservation (歸還 available)
func (s *StockCache) CancelReservation(ctx context.Context, productID int64, quantity int32) error {
	script := `
		local availKey = KEYS[1]
		local reservKey = KEYS[2]
		local quantity = tonumber(ARGV[1])
		
		local reserved = tonumber(redis.call('GET', reservKey) or 0)
		
		if reserved < quantity then
			return 0
		end
		
		redis.call('INCRBY', availKey, quantity)
		redis.call('DECRBY', reservKey, quantity)
		return 1
	`

	result, err := s.client.Eval(ctx, script, []string{
		s.availableKey(productID),
		s.reservedKey(productID),
	}, quantity).Int()

	if err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	if result == 0 {
		return fmt.Errorf("insufficient reserved stock")
	}

	return nil
}

// DeleteStock removes stock from cache
func (s *StockCache) DeleteStock(ctx context.Context, productID int64) error {
	pipe := s.client.Pipeline()
	pipe.Del(ctx, s.availableKey(productID))
	pipe.Del(ctx, s.reservedKey(productID))

	_, err := pipe.Exec(ctx)
	return err
}

// RefreshTTL refreshes the TTL of stock keys
func (s *StockCache) RefreshTTL(ctx context.Context, productID int64) error {
	pipe := s.client.Pipeline()
	pipe.Expire(ctx, s.availableKey(productID), s.ttl)
	pipe.Expire(ctx, s.reservedKey(productID), s.ttl)

	_, err := pipe.Exec(ctx)
	return err
}
