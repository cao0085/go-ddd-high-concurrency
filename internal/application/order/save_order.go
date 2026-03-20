package order

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"flash-sale-order-system/pkg/rabbitmq"
)

// SaveOrderHandler persists an order received from RabbitMQ into PostgreSQL.
type SaveOrderHandler struct {
	db *sql.DB
}

func NewSaveOrderHandler(db *sql.DB) *SaveOrderHandler {
	return &SaveOrderHandler{db: db}
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
		`INSERT INTO orders (user_id, product_id, quantity, total_price, status) VALUES ($1, $2, $3, $4, 'pending')`,
		userID, productID, msg.Quantity, totalPrice,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}
