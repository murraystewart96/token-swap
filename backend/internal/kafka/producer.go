package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/pkg/tracing"
	"go.opentelemetry.io/otel/propagation"
)

type IProducer interface {
	Produce(ctx context.Context, topic string, key, value []byte) error
	Close() error
}

type Producer struct {
	client *kafka.Producer
}

func NewProducer(cfg *config.KafkaProducer) (*Producer, error) {
	client, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"acks":              cfg.Acks,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create producer: %w", err)
	}

	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, topic string, key, value []byte) error {
	// Start a span for Kafka produce operation
	ctx, span := tracing.StartSpan(ctx, "kafka.produce")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		tracing.KafkaAttributes(topic, 0, -1)..., // partition and offset unknown at produce time
	)

	// Create headers for trace propagation
	headers := []kafka.Header{}
	propagator := propagation.TraceContext{}
	propagator.Inject(ctx, NewHeaderCarrier(&headers))

	err := p.client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:     key,
		Value:   value,
		Headers: headers,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce event: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	p.client.Close()
	return nil
}
