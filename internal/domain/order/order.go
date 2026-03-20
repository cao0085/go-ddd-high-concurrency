package order

import "time"

// Status represents the lifecycle state of an order.
type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

// Order is the aggregate root for the order domain.
type Order struct {
	ID        int64
	UserID    int64
	ProductID int64
	Quantity  int32
	Price     float64
	Status    Status
	CreatedAt time.Time
}

// NewOrder creates a new pending Order with the given parameters.
// The ID must be generated externally (e.g. via Snowflake idgen).
func NewOrder(id, userID, productID int64, quantity int32, price float64) *Order {
	return &Order{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

// TotalPrice returns the total price of the order.
func (o *Order) TotalPrice() float64 {
	return float64(o.Quantity) * o.Price
}
