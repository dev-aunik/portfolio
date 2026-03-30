// Package messaging — Kafka producer (audit events via Redpanda).
package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaBus implements ports.MessageBus using Kafka / Redpanda.
type KafkaBus struct {
	writer *kafka.Writer
}

// NewKafka creates a Kafka writer connected to the given brokers.
func NewKafka(brokers []string, topicAudit string) (*KafkaBus, error) {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topicAudit,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           5 * time.Second,
		ReadTimeout:            5 * time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}
	return &KafkaBus{writer: w}, nil
}

// Publish sends a raw JSON payload to Kafka with the topic as the message key.
func (k *KafkaBus) Publish(ctx context.Context, topic string, payload []byte) error {
	msg := kafka.Message{
		Key:   []byte(topic),
		Value: payload,
		Time:  time.Now().UTC(),
	}
	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka: write message: %w", err)
	}
	return nil
}

// Close shuts down the Kafka writer.
func (k *KafkaBus) Close() error {
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}
