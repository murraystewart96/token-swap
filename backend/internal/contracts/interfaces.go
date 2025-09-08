package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// PoolContract defines the interface for interacting with the pool contract.
// This interface contains only the methods we actually use from the generated Pool contract.
type PoolContract interface {
	// ParseSwap parses a Swap event from a blockchain log
	ParseSwap(log types.Log) (*PoolSwap, error)
	
	// ParseSync parses a Sync event from a blockchain log  
	ParseSync(log types.Log) (*PoolSync, error)
	
	// GetReserves returns the current token reserves from the pool
	GetReserves(opts *bind.CallOpts) (struct {
		MeTokenReserve  *big.Int
		YouTokenReserve *big.Int
	}, error)
}

// Ensure that the generated Pool contract implements our interface
var _ PoolContract = (*Pool)(nil)