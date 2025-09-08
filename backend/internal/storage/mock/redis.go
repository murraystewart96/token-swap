package mock

import (
	"context"

	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/stretchr/testify/mock"
)

type PoolCache struct {
	mock.Mock
}

func (m *PoolCache) SetPrice(ctx context.Context, pair string, price string) error {
	args := m.Called(ctx, pair, price)
	return args.Error(0)
}

func (m *PoolCache) GetPrice(ctx context.Context, pair string) (string, error) {
	args := m.Called(ctx, pair)
	return args.String(0), args.Error(1)
}

func (m *PoolCache) SetReserves(ctx context.Context, poolAddr string, reserves *models.PoolReserves) error {
	args := m.Called(ctx, poolAddr, reserves)
	return args.Error(0)
}

func (m *PoolCache) GetReserves(ctx context.Context, poolAddr string) (*models.PoolReserves, error) {
	args := m.Called(ctx, poolAddr)
	return args.Get(0).(*models.PoolReserves), args.Error(1)
}

func (m *PoolCache) Reset() error {
	return nil
}
