package order

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	redisinfra "flash-sale-order-system/internal/Infrastructure/persistence/redis"
	"flash-sale-order-system/pkg/rabbitmq"
)

// SaveOrderHandler persists an order received from RabbitMQ into PostgreSQL.
type SaveOrderHandler struct {
	db         *sql.DB
	stockCache *redisinfra.StockCache
}

func NewSaveOrderHandler(db *sql.DB, stockCache *redisinfra.StockCache) *SaveOrderHandler {
	return &SaveOrderHandler{db: db, stockCache: stockCache}
}

func (h *SaveOrderHandler) Handle(ctx context.Context, msg rabbitmq.OrderMessage) error {
	productID, err := strconv.ParseInt(msg.ProductID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid product_id %q: %w", msg.ProductID, err)
	}

	userID, err := strconv.ParseInt(msg.UserID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user_id %q: %w", msg.UserID, err)
	}

	totalPrice := float64(msg.Quantity) * msg.Price

	_, err = h.db.ExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, product_id, quantity, total_price, status)
		 VALUES ($1, $2, $3, $4, $5, 'pending')`,
		msg.OrderID, userID, productID, msg.Quantity, totalPrice,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := h.stockCache.ConfirmReservation(ctx, productID, int32(msg.Quantity)); err != nil {
		log.Printf("save_order: confirm reservation failed for product %d: %v", productID, err)
	}

	if _, err := h.db.ExecContext(ctx,
		`UPDATE products SET available_stock = available_stock - $1, updated_at = NOW() WHERE id = $2`,
		msg.Quantity, productID,
	); err != nil {
		log.Printf("save_order: failed to decrement product stock for product %d: %v", productID, err)
	}

	return nil
}
