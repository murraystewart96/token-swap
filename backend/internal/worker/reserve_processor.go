package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/pkg/tracing"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	MET_YOU_PAIR = "MET_YOU"
)

func (w *Worker) handleReserveEvent(ctx context.Context, key, value []byte) error {
	// Start a span for reserve event processing
	ctx, span := tracing.StartSpan(ctx, "worker.handleReserveEvent")
	defer span.End()
	var reserveEvent models.ReserveEvent
	if err := json.Unmarshal(value, &reserveEvent); err != nil {
		return fmt.Errorf("failed to unmarshal reserve event: %w", err)
	}

	// Add reserve-specific attributes to the span
	span.SetAttributes(
		tracing.BlockchainAttributes(reserveEvent.BlockNumber, reserveEvent.TxHash)...,
	)

	reserves := &models.PoolReserves{
		METAmount: reserveEvent.METReserve,
		YOUAmount: reserveEvent.YOUReserve,
	}

	// Cache current price calculated from latest reserves
	currentPrice, err := calculateMetToYouPrice(reserves)
	if err != nil {
		return fmt.Errorf("failed to calculate trading price: %w", err)
	}

	log.Info().Str("MET:YOU", currentPrice).Msg("caching latest price")

	// Start span for price caching
	ctx, priceSpan := tracing.StartSpan(ctx, "cache.SetPrice")
	err = w.poolCache.SetPrice(ctx, MET_YOU_PAIR, currentPrice)
	priceSpan.End()
	if err != nil {
		return fmt.Errorf("failed to update cache with trade price: %w", err)
	}

	log.Info().Str("MET", reserveEvent.METReserve).Str("YOU", reserves.YOUAmount).Msg("caching latest reserves")

	// Start span for reserves caching
	ctx, reservesSpan := tracing.StartSpan(ctx, "cache.SetReserves")
	err = w.poolCache.SetReserves(ctx, reserveEvent.PoolAddress, reserves)
	reservesSpan.End()
	if err != nil {
		return fmt.Errorf("failed to update reserves cache: %w", err)
	}

	log.Info().Msg("storing reserve in database")

	// Start span for database operation
	ctx, dbSpan := tracing.StartSpan(ctx, "db.CreateReserve")
	defer dbSpan.End()
	
	err = w.db.CreateReserve(&reserveEvent)
	if err != nil {
		return fmt.Errorf("failed to store reserve event in database: %w", err)
	}

	return nil
}

func calculateMetToYouPrice(reserves *models.PoolReserves) (string, error) {
	var metAmount, youAmount decimal.Decimal
	var err error

	metAmount, err = decimal.NewFromString(reserves.METAmount)
	if err != nil {
		return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", reserves.METAmount, err)
	}
	youAmount, err = decimal.NewFromString(reserves.YOUAmount)
	if err != nil {
		return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", reserves.YOUAmount, err)
	}

	// Always calculate MET:YOU (YOU per MET)
	if metAmount.IsZero() {
		return "", fmt.Errorf("MET amount is 0: can't divide by 0")
	}

	price := youAmount.Div(metAmount)
	return price.StringFixed(6), nil
}
