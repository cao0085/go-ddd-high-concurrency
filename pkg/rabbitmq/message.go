package rabbitmq

// OrderMessage is the message payload published to the orders queue.
type OrderMessage struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
