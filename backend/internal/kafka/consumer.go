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
	StartConsuming(ctx context.Context, handlers EventHandlers) error
	Close() error
}

type Consumer struct {
	client *kafka.Consumer
}

type EventHandler func(key, value []byte) error
type EventHandlers map[string]EventHandler

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

func (c *Consumer) StartConsuming(ctx context.Context, topicHandlers EventHandlers) error {
	topics := make([]string, 0, len(topicHandlers))
	for topic := range topicHandlers {
		topics = append(topics, topic)
	}

	log.Info().Msgf("subscribing to topics: %v", topics)

	err := c.client.SubscribeTopics(topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics (%v): %w", topics, err)
	}

	log.Info().Msgf("consuming...")

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
					log.Error().Err(err).Msg("failed to read from topic")
				}
				continue
			}

			if handler, ok := topicHandlers[*message.TopicPartition.Topic]; ok {
				if err := handler(message.Key, message.Value); err != nil {
					log.Error().Err(err).Msg("message handler failed")
				}
			} else {
				log.Warn().Msgf("no handler for topic: %s", *message.TopicPartition.Topic)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.client.Close()
}
