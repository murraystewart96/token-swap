package mock

import (
	"time"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/stretchr/testify/mock"
)

type DB struct {
	mock.Mock
}

// Trade operations
func (m *DB) CreateTrade(trade *models.TradeEvent) error {
	args := m.Called(trade)
	return args.Error(0)
}

func (m *DB) GetTradesByTimeRange(start, end time.Time) ([]*models.TradeEvent, error) {
	args := m.Called(start, end)
	return args.Get(0).([]*models.TradeEvent), args.Error(1)
}

func (m *DB) GetTradesByCursor(cursorBlock uint64, cursorTx uint, limit int) ([]*models.TradeEvent, error) {
	args := m.Called(cursorBlock, cursorTx, limit)
	return args.Get(0).([]*models.TradeEvent), args.Error(1)
}

func (m *DB) UpdateConfirmedTrades(confirmationThreshold uint64) error {
	args := m.Called(confirmationThreshold)
	return args.Error(0)
}

// Reserve operations
func (m *DB) CreateReserve(reserve *models.ReserveEvent) error {
	args := m.Called(reserve)
	return args.Error(0)
}

func (m *DB) GetReservesByTimeRange(start, end time.Time) ([]*models.ReserveEvent, error) {
	args := m.Called(start, end)
	return args.Get(0).([]*models.ReserveEvent), args.Error(1)
}

func (m *DB) UpdateConfirmedReserves(confirmationThreshold uint64) error {
	args := m.Called(confirmationThreshold)
	return args.Error(0)
}

// Analytics
func (m *DB) GetVolumeAnalytics(start, end time.Time, token string) (*models.VolumeResponse, error) {
	args := m.Called(start, end, token)
	return args.Get(0).(*models.VolumeResponse), args.Error(1)
}

func (m *DB) GetPriceHistory(start, end time.Time, interval time.Duration) (*models.PriceHistoryResponse, error) {
	args := m.Called(start, end, interval)
	return args.Get(0).(*models.PriceHistoryResponse), args.Error(1)
}

func (m *DB) GetActivityAnalytics(start, end time.Time) (*models.ActivityResponse, error) {
	args := m.Called(start, end)
	return args.Get(0).(*models.ActivityResponse), args.Error(1)
}

func (m *DB) RollbackEvents(blockNumber uint64) error {
	args := m.Called(blockNumber)
	return args.Error(0)
}

func (m *DB) Close() {
	m.Called()
}

func (m *DB) Exec(query string) error {
	return nil
}
