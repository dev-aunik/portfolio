// Package messaging provides RabbitMQ and Kafka implementations of ports.MessageBus.
package messaging

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQBus implements ports.MessageBus using RabbitMQ.
type RabbitMQBus struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchange     string
	queueContact string
}

// NewRabbitMQ connects to RabbitMQ and declares the exchange and queues.
func NewRabbitMQ(url, exchange, queueContact string) (*RabbitMQBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: open channel: %w", err)
	}

	// Declare a durable topic exchange
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: declare exchange: %w", err)
	}

	// Declare the contact submissions queue (durable, persistent)
	q, err := ch.QueueDeclare(queueContact, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: declare queue %q: %w", queueContact, err)
	}

	// Bind the queue to the exchange with routing key "contact.new"
	if err := ch.QueueBind(q.Name, "contact.new", exchange, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: bind queue: %w", err)
	}

	return &RabbitMQBus{
		conn:         conn,
		ch:           ch,
		exchange:     exchange,
		queueContact: queueContact,
	}, nil
}

// Publish sends a JSON payload to the RabbitMQ exchange with the given routing key.
func (r *RabbitMQBus) Publish(ctx context.Context, routingKey string, payload []byte) error {
	err := r.ch.PublishWithContext(ctx,
		r.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         payload,
		},
	)
	if err != nil {
		return fmt.Errorf("rabbitmq: publish to %q: %w", routingKey, err)
	}
	return nil
}

// Close gracefully shuts down the RabbitMQ connection.
func (r *RabbitMQBus) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
