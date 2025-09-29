package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IClient defines the interface for Ethereum client operations
type IClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

// Client wraps the go-ethereum ethclient to implement IClient interface
type Client struct {
	client *ethclient.Client
}

// NewClient creates a new Ethereum client instance
func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(fmt.Sprintf("ws://%s", rpcURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// BlockByNumber returns the block with the given number
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.client.BlockByNumber(ctx, number)
}

// SubscribeFilterLogs subscribes to log events matching the given filter query
func (c *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return c.client.SubscribeFilterLogs(ctx, q, ch)
}

// Close closes the underlying client connection
func (c *Client) Close() {
	c.client.Close()
}

// GetUnderlyingClient returns the underlying ethclient.Client for contract creation
func (c *Client) GetUnderlyingClient() *ethclient.Client {
	return c.client
}