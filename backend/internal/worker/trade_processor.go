package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/pkg/tracing"
	"github.com/rs/zerolog/log"
)

func (w *Worker) handleTradeEvent(ctx context.Context, key, value []byte) error {
	// Start a span for trade event processing
	ctx, span := tracing.StartSpan(ctx, "worker.handleTradeEvent")
	defer span.End()

	tradeEvent := &models.TradeEvent{}
	if err := json.Unmarshal(value, tradeEvent); err != nil {
		return fmt.Errorf("failed to unmarshal trade history event: %w", err)
	}

	// Add trade-specific attributes to the span
	span.SetAttributes(
		tracing.BlockchainAttributes(tradeEvent.BlockNumber, tradeEvent.TxHash)...,
	)

	log.Info().Msg("storing trade in database")

	// Start a span for database operation
	ctx, dbSpan := tracing.StartSpan(ctx, "db.CreateTrade")
	defer dbSpan.End()

	err := w.db.CreateTrade(tradeEvent)
	if err != nil {
		return fmt.Errorf("failed to store trade event in database: %w", err)
	}

	return nil
}
