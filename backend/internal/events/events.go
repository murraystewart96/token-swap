package events

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/contracts"
	"github.com/murraystewart96/token-swap/internal/kafka"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/rs/zerolog/log"
)

var (
	SwapEventSignature = common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")
	SyncEventSignature = common.HexToHash("0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a")
)

type EventClient struct {
	ethClient    *ethclient.Client
	contractAddr common.Address
	poolContract contracts.PoolContract
	producer     kafka.IProducer
}

func NewClient(cfg *config.Listener, producer kafka.IProducer) (*EventClient, error) {
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
		producer:     producer,
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
				err := ec.handleSwapEvent(&eventLog)
				if err != nil {
					log.Error().Err(err).Msg("swap event handler failed")
				}

			case SyncEventSignature:
				err := ec.handleSyncEvent(&eventLog)
				if err != nil {
					log.Error().Err(err).Msg("sync event handler failed")
				}
			}
		}
	}
}

func (ec *EventClient) handleSwapEvent(eventLog *types.Log) error {
	if swapEvent, err := ec.poolContract.ParseSwap(*eventLog); err == nil {
		logSwapEvent(swapEvent)

		tokenIn, tokenOut, amountIn, amountOut := determineTradeDirection(swapEvent)

		tradeEvent := models.TradeEvent{
			TxHash:           swapEvent.Raw.TxHash.Hex(),
			TransactionIndex: swapEvent.Raw.TxIndex,
			BlockNumber:      swapEvent.Raw.BlockNumber,
			Timestamp:        int64(swapEvent.Raw.BlockTimestamp),
			Sender:           swapEvent.Sender.Hex(),
			Recipient:        swapEvent.To.Hex(),
			TokenIn:          tokenIn,
			TokenOut:         tokenOut,
			AmountIn:         amountIn,
			AmountOut:        amountOut,
			PoolAddress:      swapEvent.Raw.Address.Hex(),
			EventType:        "swap",
		}

		// Convert to JSON and publish
		tradeEventJSON, err := json.Marshal(tradeEvent)
		if err != nil {
			return fmt.Errorf("failed to marshal trade event: %w", err)
		}

		err = ec.producer.Produce(config.TradeHistoryTopic, []byte(tradeEvent.TxHash), tradeEventJSON)
		if err != nil {
			return fmt.Errorf("failed to produce trade event: %w", err)
		}
	}

	return nil
}

func (ec *EventClient) handleSyncEvent(eventLog *types.Log) error {
	if syncEvent, err := ec.poolContract.ParseSync(*eventLog); err == nil {
		logSyncEvent(syncEvent)

		reserveEvent := models.ReserveEvent{
			TxHash:      syncEvent.Raw.TxHash.Hex(),
			BlockNumber: syncEvent.Raw.BlockNumber,
			METReserve:  syncEvent.MeTokenAmount.String(),
			YOUReserve:  syncEvent.YouTokenAmount.String(),
			PoolAddress: syncEvent.Raw.Address.Hex(),
			EventType:   "sync",
		}

		// Convert to JSON and publish
		reserveEventJSON, err := json.Marshal(reserveEvent)
		if err != nil {
			return fmt.Errorf("failed to marshal reserve event: %w", err)
		}

		err = ec.producer.Produce(config.ReserveHistoryTopic, []byte(reserveEvent.TxHash), reserveEventJSON)
		if err != nil {
			return fmt.Errorf("failed to produce reserve event: %w", err)
		}
	}

	return nil
}

func determineTradeDirection(swapEvent *contracts.PoolSwap) (tokenIn, tokenOut, amountIn, amountOut string) {
	if swapEvent.MeTokenIn.Cmp(big.NewInt(0)) > 0 {
		// MET → YOU trade
		return "MET", "YOU", swapEvent.MeTokenIn.String(), swapEvent.YouTokenOut.String()
	} else {
		// YOU → MET trade
		return "YOU", "MET", swapEvent.YouTokenIn.String(), swapEvent.MeTokenOut.String()
	}
}

func logSwapEvent(swapEvent *contracts.PoolSwap) {
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

func logSyncEvent(syncEvent *contracts.PoolSync) {
	log.Info().
		Str("event_type", "sync").
		Str("me_token_amount", syncEvent.MeTokenAmount.String()).
		Str("you_token_amount", syncEvent.YouTokenAmount.String()).
		Str("tx_hash", syncEvent.Raw.TxHash.Hex()).
		Uint64("block_number", syncEvent.Raw.BlockNumber).
		Msg("Sync event received")
}
