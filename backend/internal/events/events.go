package events

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/contracts"
	"github.com/rs/zerolog/log"
)

var (
	SwapEventSignature = common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")
	SyncEventSignature = common.HexToHash("0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a")
)

type EventClient struct {
	ethClient    *ethclient.Client
	contractAddr common.Address
	poolContract *contracts.Pool
}

func NewClient(cfg *config.Events) (*EventClient, error) {
	client, err := ethclient.Dial(fmt.Sprintf("ws://%s", cfg.RPCUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	contractAddr := common.HexToAddress(cfg.ContractAddr)
	poolContract, err := contracts.NewPool(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool contract interface: %w", err)
	}

	return &EventClient{
		ethClient:    client,
		contractAddr: contractAddr,
		poolContract: poolContract,
	}, nil

}

func (ec *EventClient) Listen(ctx context.Context) error {
	// One subscription for all events from this contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ec.contractAddr},
	}

	logs := make(chan types.Log)
	sub, err := ec.ethClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to contract events: %w", err)
	}

	log.Info().Msgf("listening for events on %s", ec.contractAddr.String())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return err
		case eventLog := <-logs:

			switch eventLog.Topics[0] {
			case SwapEventSignature:
				if swapEvent, err := ec.poolContract.ParseSwap(eventLog); err == nil {
					log.Info().
						Str("event_type", "swap").
						Str("sender", swapEvent.Sender.Hex()).
						Str("to", swapEvent.To.Hex()).
						Str("me_token_in", swapEvent.MeTokenIn.String()).
						Str("you_token_in", swapEvent.YouTokenIn.String()).
						Str("me_token_out", swapEvent.MeTokenOut.String()).
						Str("you_token_out", swapEvent.YouTokenOut.String()).
						Str("tx_hash", swapEvent.Raw.TxHash.Hex()).
						Uint64("block_number", swapEvent.Raw.BlockNumber).
						Msg("Swap event received")
				}
			case SyncEventSignature:
				if syncEvent, err := ec.poolContract.ParseSync(eventLog); err == nil {
					log.Info().
						Str("event_type", "sync").
						Str("me_token_amount", syncEvent.MeTokenAmount.String()).
						Str("you_token_amount", syncEvent.YouTokenAmount.String()).
						Str("tx_hash", syncEvent.Raw.TxHash.Hex()).
						Uint64("block_number", syncEvent.Raw.BlockNumber).
						Msg("Sync event received")
				}
			}
		}
	}
}
