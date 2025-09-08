package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/contracts"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/internal/storage"
	"github.com/murraystewart96/token-swap/internal/worker"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	TokenTypeME  uint8 = 0
	TokenTypeYOU uint8 = 1
)

type Sync struct {
	ethClient    *ethclient.Client
	contractAddr common.Address
	poolContract contracts.PoolContract
	db           storage.DB
	poolCache    storage.PoolCache
}

func NewSync(cfg *config.Sync, poolCache storage.PoolCache, db storage.DB) (*Sync, error) {
	client, err := ethclient.Dial(fmt.Sprintf("ws://%s", cfg.Listener.RPCUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	contractAddr := common.HexToAddress(cfg.Listener.ContractAddr)
	poolContract, err := contracts.NewPool(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool contract interface: %w", err)
	}

	return &Sync{
		ethClient:    client,
		contractAddr: contractAddr,
		poolContract: poolContract,
		poolCache:    poolCache,
		db:           db,
	}, nil
}

func (s *Sync) Start(ctx context.Context) error {
	log.Info().Msg("starting sync...")

	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			err := s.Sync(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to sync pool contract state")
			}

			err = s.updateConfirmedEvents(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to update confirmed events")
			}
		case <-ctx.Done():

			return nil
		}
	}
}

func (s *Sync) Sync(ctx context.Context) error {
	// Sync reserves first
	reserves, err := s.poolContract.GetReserves(nil)
	if err != nil {
		return fmt.Errorf("failed to get current reserves: %w", err)
	}

	// Calculate price from reserves (same as AMM formula)
	metReserve := decimal.NewFromBigInt(reserves.MeTokenReserve, 0)
	youReserve := decimal.NewFromBigInt(reserves.YouTokenReserve, 0)

	if metReserve.IsZero() {
		return fmt.Errorf("MET reserves are empty: can't divide by zero")
	}

	// Cache Current price
	currentPrice := youReserve.Div(metReserve) // YOU per MET

	log.Info().Str("MET:YOU", currentPrice.StringFixed(6)).Msg("syncing cache with latest price")

	err = s.poolCache.SetPrice(ctx, worker.MET_YOU_PAIR, currentPrice.StringFixed(6))
	if err != nil {
		return fmt.Errorf("failed to sync cache with trade price: %w", err)
	}

	// Cache reserves
	poolReserves := &models.PoolReserves{
		METAmount: reserves.MeTokenReserve.String(),
		YOUAmount: reserves.YouTokenReserve.String(),
	}

	log.Info().Str("MET", reserves.MeTokenReserve.String()).Str("YOU", reserves.YouTokenReserve.String()).Msg("syncing cache with latest reserves")

	err = s.poolCache.SetReserves(ctx, s.contractAddr.Hex(), poolReserves)
	if err != nil {
		return fmt.Errorf("failed to update reserves cache: %w", err)
	}

	return nil
}

func (s *Sync) updateConfirmedEvents(ctx context.Context) error {
	currentBlock, err := s.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	confirmationThreshold := currentBlock - 12

	if err := s.db.UpdateConfirmedTrades(confirmationThreshold); err != nil {
		return fmt.Errorf("failed to update confirmed trades: %w", err)
	}

	if err := s.db.UpdateConfirmedReserves(confirmationThreshold); err != nil {
		return fmt.Errorf("failed to update confirmed reserves: %w", err)
	}

	return nil
}
