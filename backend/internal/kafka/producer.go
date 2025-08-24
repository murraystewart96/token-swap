package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/murraystewart96/token-swap/internal/config"
)

type IProducer interface {
	Produce(topic string, key, value []byte) error
}

type Producer struct {
	client *kafka.Producer
}

func NewProducer(cfg *config.KafkaProducer) (*Producer, error) {
	client, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"acks":              cfg.Acks})
	if err != nil {
		return nil, fmt.Errorf("Failed to create producer: %w", err)
	}

	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) Produce(topic string, key, value []byte) error {
	err := p.client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: []byte(value),
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce event: %w", err)
	}

	return nil
}
