package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	MET_YOU_PAIR = "MET_YOU"
)

func (w *Worker) handleTradeEvent(key, value []byte) error {
	tradeEvent := &models.TradeEvent{}
	if err := json.Unmarshal(value, tradeEvent); err != nil {
		return fmt.Errorf("failed to unmarshal trade history event: %w", err)
	}

	currentPrice, err := calculateMetToYouPrice(tradeEvent)
	if err != nil {
		return fmt.Errorf("failed to calculate trading price: %w", err)
	}

	log.Info().Str("MET:YOU", currentPrice).Msg("caching latest price")

	err = w.poolCache.SetPrice(context.Background(), MET_YOU_PAIR, currentPrice)
	if err != nil {
		return fmt.Errorf("failed to update cache with trade price: %w", err)
	}

	log.Info().Msg("storing price in database")

	err = w.db.CreateTrade(tradeEvent)
	if err != nil {
		return fmt.Errorf("failed to store trade event in database: %w", err)
	}

	return nil
}

func calculateMetToYouPrice(tradeEvent *models.TradeEvent) (string, error) {
	var metAmount, youAmount decimal.Decimal
	var err error

	if tradeEvent.TokenIn == "MET" {
		// MET → YOU trade
		metAmount, err = decimal.NewFromString(tradeEvent.AmountIn)
		if err != nil {
			return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", tradeEvent.AmountIn, err)
		}
		youAmount, err = decimal.NewFromString(tradeEvent.AmountOut)
		if err != nil {
			return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", tradeEvent.AmountOut, err)
		}
	} else {
		// YOU → MET trade - need to flip the calculation
		metAmount, err = decimal.NewFromString(tradeEvent.AmountOut)
		if err != nil {
			return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", tradeEvent.AmountOut, err)
		}
		youAmount, err = decimal.NewFromString(tradeEvent.AmountIn)
		if err != nil {
			return "", fmt.Errorf("failed to convert string (%s) to decimal: %w", tradeEvent.AmountIn, err)
		}
	}

	// Always calculate MET:YOU (YOU per MET)
	if metAmount.IsZero() {
		return "", fmt.Errorf("MET amount is 0: can't divide by 0")
	}

	price := youAmount.Div(metAmount)
	return price.StringFixed(6), nil
}
