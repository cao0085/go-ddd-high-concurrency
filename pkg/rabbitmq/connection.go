package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const OrderQueue = "orders"

// NewConnection creates an AMQP connection and channel.
// Callers are responsible for closing both when done.
func NewConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// Declare the durable orders queue so both producer and consumer share the same definition.
	_, err = ch.QueueDeclare(
		OrderQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
