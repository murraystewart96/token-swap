package worker

import (
	"context"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/pkg/kafka"
	"github.com/murraystewart96/token-swap/internal/storage"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	consumer      kafka.IConsumer
	db            storage.DB
	poolCache     storage.PoolCache
	eventHandlers kafka.EventHandlers
}

func New(consumer kafka.IConsumer, topics []string, poolCache storage.PoolCache, db storage.DB) (*Worker, error) {
	worker := &Worker{
		consumer:  consumer,
		poolCache: poolCache,
		db:        db,
	}

	eventHandlers := kafka.EventHandlers{
		config.TradeHistoryTopic:   worker.handleTradeEvent,
		config.ReserveHistoryTopic: worker.handleReserveEvent,
	}

	// Assign configured topic handlers
	activeHandlers := make(kafka.EventHandlers)
	for _, topic := range topics {
		handler, found := eventHandlers[topic]
		if !found {
			return nil, fmt.Errorf("no event handler for topic: %s", topic)
		}
		activeHandlers[topic] = handler
	}

	worker.eventHandlers = activeHandlers

	return worker, nil
}

func (w *Worker) Start(ctx context.Context) error {
	log.Info().Msg("starting worker process")

	if err := w.processEvents(ctx); err != nil {
		return fmt.Errorf("failed to process events: %w", err)
	}

	return nil
}

func (w *Worker) processEvents(ctx context.Context) error {
	log.Info().Msgf("processing events...")

	err := w.consumer.StartConsuming(ctx, w.eventHandlers)
	if err != nil {
		return fmt.Errorf("failed to start consuming topic: %w", err)
	}

	return nil
}

func (w *Worker) SetEventHandlers(handlers kafka.EventHandlers) {
	w.eventHandlers = handlers
}
