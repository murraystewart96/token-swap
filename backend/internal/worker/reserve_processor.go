package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/rs/zerolog/log"
)

func (w *Worker) handleReserveEvent(key, value []byte) error {
	var reserveEvent models.ReserveEvent
	if err := json.Unmarshal(value, &reserveEvent); err != nil {
		return fmt.Errorf("failed to unmarshal reserve event: %w", err)
	}

	reserves := &models.PoolReserves{
		METAmount: reserveEvent.METReserve,
		YOUAmount: reserveEvent.YOUReserve,
	}

	log.Info().Str("MET", reserveEvent.METReserve).Str("YOU", reserves.YOUAmount).Msg("caching latest reserves")

	err := w.poolCache.SetReserves(context.Background(), reserveEvent.PoolAddress, reserves)
	if err != nil {
		return fmt.Errorf("failed to update reserves cache: %w", err)
	}

	log.Info().Msg("storing reserve in database")

	err = w.db.CreateReserve(&reserveEvent)
	if err != nil {
		return fmt.Errorf("failed to store reserve event in database: %w", err)
	}

	return nil
}
