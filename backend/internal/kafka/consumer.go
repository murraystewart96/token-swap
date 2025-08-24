package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/rs/zerolog/log"
)

type IConsumer interface {
	StartConsuming(ctx context.Context, topic string, handler MessageHandler) error
}

type Consumer struct {
	client *kafka.Consumer
	topic  string
}

type MessageHandler func(key, value []byte) error

func NewConsumer(cfg *config.KafkaConsumer) (*Consumer, error) {
	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.OffsetReset})
	if err != nil {
		return nil, fmt.Errorf("Failed to create consumer: %w", err)
	}

	return &Consumer{
		client: client,
	}, nil
}

func (c *Consumer) StartConsuming(ctx context.Context, topic string, handler MessageHandler) error {
	err := c.client.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic (%s): %w", topic, err)
	}

	for {
		select {
		case <-ctx.Done():
			c.client.Close()

			return nil
		default:
			message, err := c.client.ReadMessage(300 * time.Millisecond)
			if err != nil {
				// Don't log timeouts (normal polling behaviour)
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					log.Error().Err(err).Msgf("failed to read from topic: %s", topic)
				}
				continue
			}

			if err := handler(message.Key, message.Value); err != nil {
				log.Error().Err(err).Msg("message handler failed")
			}
		}
	}
}
