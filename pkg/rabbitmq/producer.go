package rabbitmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderProducer publishes order messages to the orders queue.
type OrderProducer struct {
	ch *amqp.Channel
}

func NewOrderProducer(ch *amqp.Channel) *OrderProducer {
	return &OrderProducer{ch: ch}
}

func (p *OrderProducer) Publish(msg OrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.ch.Publish(
		"",         // default exchange
		OrderQueue, // routing key = queue name
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
