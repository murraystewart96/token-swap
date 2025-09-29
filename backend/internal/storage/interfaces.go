package storage

import (
	"context"
	"time"

	"github.com/murraystewart96/token-swap/internal/models"
)

type PoolCache interface {
	SetPrice(ctx context.Context, pair string, price string) error
	GetPrice(ctx context.Context, pair string) (string, error)
	SetReserves(ctx context.Context, poolAddr string, reserves *models.PoolReserves) error
	GetReserves(ctx context.Context, poolAddr string) (*models.PoolReserves, error)
	Reset() error
}

type DB interface {
	// Trade operations
	CreateTrade(trade *models.TradeEvent) error
	GetTradesByTimeRange(start, end time.Time) ([]*models.TradeEvent, error)
	GetTradesByCursor(cursorBlock uint64, cursorTx uint, limit int) ([]*models.TradeEvent, error)
	UpdateConfirmedTrades(confirmationThreshold uint64) error

	// Reserve operations
	CreateReserve(reserve *models.ReserveEvent) error
	GetReservesByTimeRange(start, end time.Time) ([]*models.ReserveEvent, error)
	UpdateConfirmedReserves(confirmationThreshold uint64) error

	// Analytics
	GetVolumeAnalytics(start, end time.Time, token string) (*models.VolumeResponse, error)
	GetPriceHistory(start, end time.Time, interval time.Duration) (*models.PriceHistoryResponse, error)
	GetActivityAnalytics(start, end time.Time) (*models.ActivityResponse, error)

	RollbackEvents(blockNumber uint64) error

	// Infrastructure operations
	Exec(query string) error
	Close()
}
