package order

import (
	"context"
	"fmt"

	redisinfra "flash-sale-order-system/internal/Infrastructure/persistence/redis"
	"flash-sale-order-system/pkg/rabbitmq"
)

type PlaceOrderCommand struct {
	OrderID   string
	UserID    string
	ProductID int64
	Quantity  int32
	Price     float64
}

type PlaceOrderHandler struct {
	stockCache *redisinfra.StockCache
	producer   *rabbitmq.OrderProducer
}

func NewPlaceOrderHandler(stockCache *redisinfra.StockCache, producer *rabbitmq.OrderProducer) *PlaceOrderHandler {
	return &PlaceOrderHandler{
		stockCache: stockCache,
		producer:   producer,
	}
}

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd PlaceOrderCommand) error {
	// 1. Atomically reserve stock in Redis
	ok, err := h.stockCache.Reserve(ctx, cmd.ProductID, cmd.Quantity)
	if err != nil {
		return fmt.Errorf("stock reserve error: %w", err)
	}
	if !ok {
		return fmt.Errorf("insufficient stock for product %d", cmd.ProductID)
	}

	// 2. Publish order message to RabbitMQ
	msg := rabbitmq.OrderMessage{
		OrderID:   cmd.OrderID,
		UserID:    cmd.UserID,
		ProductID: fmt.Sprintf("%d", cmd.ProductID),
		Quantity:  int(cmd.Quantity),
		Price:     cmd.Price,
	}
	if err := h.producer.Publish(msg); err != nil {
		// Rollback stock reservation on publish failure
		_ = h.stockCache.CancelReservation(ctx, cmd.ProductID, cmd.Quantity)
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}
