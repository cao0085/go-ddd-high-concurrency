package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderConsumer reads order messages from the orders queue.
type OrderConsumer struct {
	ch *amqp.Channel
}

func NewOrderConsumer(ch *amqp.Channel) *OrderConsumer {
	return &OrderConsumer{ch: ch}
}

// Start begins consuming messages and calls handler for each one.
// It blocks until ctx is cancelled or the channel is closed.
func (c *OrderConsumer) Start(ctx context.Context, handler func(OrderMessage) error) error {
	msgs, err := c.ch.Consume(
		OrderQueue,
		"",    // consumer tag (auto)
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				var msg OrderMessage
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					log.Printf("rabbitmq: failed to unmarshal order message: %v", err)
					d.Nack(false, false)
					continue
				}
				if err := handler(msg); err != nil {
					log.Printf("rabbitmq: handler error for order %s: %v", msg.OrderID, err)
					d.Nack(false, true) // requeue
					continue
				}
				d.Ack(false)
			}
		}
	}()

	return nil
}
