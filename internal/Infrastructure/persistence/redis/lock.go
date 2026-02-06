package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DistributedLock provides distributed locking using Redis
type DistributedLock struct {
	client *redis.Client
}

// NewDistributedLock creates a new DistributedLock instance
func NewDistributedLock(client *redis.Client) *DistributedLock {
	return &DistributedLock{
		client: client,
	}
}

// lockKey generates Redis key for lock
func (l *DistributedLock) lockKey(resource string) string {
	return fmt.Sprintf("lock:%s", resource)
}

// Acquire attempts to acquire a lock
// Returns true if lock acquired, false otherwise
func (l *DistributedLock) Acquire(ctx context.Context, resource string, ttl time.Duration) (bool, error) {
	key := l.lockKey(resource)

	// SET key value NX EX ttl
	// NX: only set if not exists
	// EX: set expiration time
	result, err := l.client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

// Release releases a lock
func (l *DistributedLock) Release(ctx context.Context, resource string) error {
	key := l.lockKey(resource)

	err := l.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return nil
}

// AcquireWithRetry attempts to acquire a lock with retry
func (l *DistributedLock) AcquireWithRetry(ctx context.Context, resource string, ttl time.Duration, maxRetries int, retryInterval time.Duration) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		acquired, err := l.Acquire(ctx, resource, ttl)
		if err != nil {
			return false, err
		}

		if acquired {
			return true, nil
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryInterval):
			continue
		}
	}

	return false, nil
}

// WithLock executes a function with a distributed lock
func (l *DistributedLock) WithLock(ctx context.Context, resource string, ttl time.Duration, fn func() error) error {
	acquired, err := l.Acquire(ctx, resource, ttl)
	if err != nil {
		return err
	}

	if !acquired {
		return fmt.Errorf("failed to acquire lock for resource: %s", resource)
	}

	defer l.Release(ctx, resource)

	return fn()
}
