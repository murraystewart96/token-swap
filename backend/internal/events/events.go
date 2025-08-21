package events

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
)

type EventClient struct {
	ethClient    *ethclient.Client
	contractAddr common.Address
}

func NewClient(cfg config.Events) (*EventClient, error) {
	client, err := ethclient.Dial(fmt.Sprintf("ws://%s", cfg.RPCUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}
	return &EventClient{
		ethClient:    client,
		contractAddr: common.HexToAddress(cfg.ContractAddr),
	}, nil

}

func (ec *EventClient) Listen() error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ec.contractAddr},
	}

	logs := make(chan types.Log)
	sub, err := ec.ethClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs:

			}
		}
	}()

	return nil
}
