package order

import (
	"context"
	"fmt"
	"strconv"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	redisinfra "flash-sale-order-system/internal/Infrastructure/persistence/redis"
	domain "flash-sale-order-system/internal/domain/order"
	"flash-sale-order-system/pkg/rabbitmq"
)

type PlaceOrderCommand struct {
	UserID    string
	ProductID int64
	Quantity  int32
	Price     float64
}

type PlaceOrderHandler struct {
	stockCache *redisinfra.StockCache
	producer   *rabbitmq.OrderProducer
	idGen      *idgen.IDGenerator
}

func NewPlaceOrderHandler(stockCache *redisinfra.StockCache, producer *rabbitmq.OrderProducer, idGen *idgen.IDGenerator) *PlaceOrderHandler {
	return &PlaceOrderHandler{
		stockCache: stockCache,
		producer:   producer,
		idGen:      idGen,
	}
}

// Handle reserves stock, creates an Order aggregate, and publishes to RabbitMQ.
// Returns the generated order ID on success.
func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd PlaceOrderCommand) (int64, error) {
	// 1. Atomically reserve stock in Redis
	ok, err := h.stockCache.Reserve(ctx, cmd.ProductID, cmd.Quantity)
	if err != nil {
		return 0, fmt.Errorf("stock reserve error: %w", err)
	}
	if !ok {
		return 0, fmt.Errorf("insufficient stock for product %d", cmd.ProductID)
	}

	// 2. Build Order aggregate (ID generated internally)
	userID, err := strconv.ParseInt(cmd.UserID, 10, 64)
	if err != nil {
		_ = h.stockCache.CancelReservation(ctx, cmd.ProductID, cmd.Quantity)
		return 0, fmt.Errorf("invalid user_id %q: %w", cmd.UserID, err)
	}

	order := domain.NewOrder(h.idGen.Generate(), userID, cmd.ProductID, cmd.Quantity, cmd.Price)

	// 3. Publish order message to RabbitMQ
	msg := rabbitmq.OrderMessage{
		OrderID:   order.ID,
		UserID:    cmd.UserID,
		ProductID: fmt.Sprintf("%d", order.ProductID),
		Quantity:  int(order.Quantity),
		Price:     order.Price,
	}
	if err := h.producer.Publish(msg); err != nil {
		// Rollback stock reservation on publish failure
		_ = h.stockCache.CancelReservation(ctx, cmd.ProductID, cmd.Quantity)
		return 0, fmt.Errorf("failed to publish order: %w", err)
	}

	return order.ID, nil
}
