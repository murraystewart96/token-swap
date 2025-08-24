package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/kafka"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	tradeConsumer kafka.IConsumer // Consumes trade-history topic
	// syncConsumer  *kafka.Consumer  // Consumes pool-stats topic
	// db           *database.DB
	// cache        *redis.Cache
}

func New(tradeConsumer kafka.IConsumer) *Worker {
	return &Worker{
		tradeConsumer: tradeConsumer,
	}
}

func (w *Worker) Start(ctx context.Context) {
	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return w.processTradeHistoryEvents(ctx)
	})

	// errGroup.Go(func() error {
	// 	return w.processTradeHistoryEvents(ctx)
	// })

	errGroup.Wait()
}

func (w *Worker) processTradeHistoryEvents(ctx context.Context) error {
	err := w.tradeConsumer.StartConsuming(ctx, config.TradeHistoryTopic, handleTradeHistoryEvent)
	if err != nil {
		return fmt.Errorf("failed to start consuming topic (%s): %w", config.TradeHistoryTopic, err)
	}

	return nil
}

func handleTradeHistoryEvent(key, value []byte) error {
	tradeEvent := &models.TradeEvent{}
	if err := json.Unmarshal(value, tradeEvent); err != nil {
		log.Error().Err(err).
			Str("message_key", string(key)).
			Str("raw_value", string(value)).
			Msg("failed to unmarshal trade history event")

		return fmt.Errorf("failed to unmarshal trade history event")
	}

	logTradeHistoryEvent(string(key), tradeEvent)

	return nil
}

func logTradeHistoryEvent(key string, tradeEvent *models.TradeEvent) {
	log.Info().
		Str("message_key", key).
		Str("tx_hash", tradeEvent.TxHash).
		Uint64("block_number", tradeEvent.BlockNumber).
		Str("sender", tradeEvent.Sender).
		Str("trade_direction", fmt.Sprintf("%s → %s", tradeEvent.TokenIn, tradeEvent.TokenOut)).
		Str("amounts", fmt.Sprintf("%s %s for %s %s", tradeEvent.AmountIn, tradeEvent.TokenIn, tradeEvent.AmountOut, tradeEvent.TokenOut)).
		Str("pool_address", tradeEvent.PoolAddress).
		Msg("Trade history event processed")
}
