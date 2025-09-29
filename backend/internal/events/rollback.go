package events

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (ec *EventClient) CheckForChainReorg(ctx context.Context, event *types.Log) (bool, error) {
	// Check if we've seen this block number before
	if existingHash, seen := ec.recentBlocks[event.BlockNumber]; seen {
		if existingHash != event.BlockHash {
			// Conflict detected - query the canonical chain
			canonicalHash, err := ec.getCanonicalBlockHash(ctx, event.BlockNumber)
			if err != nil {
				return false, fmt.Errorf("failed to get canonical hash: %w", err)
			}

			switch canonicalHash {
			case event.BlockHash:
				// New event is on canonical chain, old one was reorged
				err := ec.rollbackFromBlock(event.BlockNumber)
				if err != nil {
					return true, err
				}
				ec.addBlock(event.BlockNumber, event.BlockHash)

				return true, nil
			case existingHash:
				// Existing hash is canoniacal, ignore this event
				return false, nil
			default:
				// Neither hash matches canonical - both were reorged
				err := ec.rollbackFromBlock(event.BlockNumber)
				if err != nil {
					return true, err
				}
				ec.addBlock(event.BlockNumber, canonicalHash)

				return true, nil
			}
		}
	}

	// Mark block as seen with given hash
	ec.addBlock(event.BlockNumber, event.BlockHash)

	return false, nil
}

// rollbackFromBlock removes events from potentially false chain.
// NOTE: Current implementation has limitation with out-of-order event delivery
// during reorgs - may delete valid canonical events if higher block numbers
// arrive first. Production system would need to track canonical status.
func (ec *EventClient) rollbackFromBlock(blockNumber uint64) error {
	// Delete all events >= this block number (they're on false chain)
	err := ec.db.RollbackEvents(blockNumber)
	if err != nil {
		return err
	}

	// Also clear invalid blocks from cache
	for blockNum := range ec.recentBlocks {
		if blockNum >= blockNumber {
			delete(ec.recentBlocks, blockNum)
		}
	}

	return nil
}

func (ec *EventClient) getCanonicalBlockHash(ctx context.Context, blockNum uint64) (common.Hash, error) {
	block, err := ec.ethClient.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
	if err != nil {
		return common.Hash{}, err
	}
	return block.Hash(), nil
}

// Adds block to recent blocks (keeps recent blocks lenght to MAX_RECENT_BLOCKS)
func (ec *EventClient) addBlock(num uint64, hash common.Hash) {
	ec.recentBlocks[num] = hash

	// Remove oldest block
	// In production we would use LRU cache instead of this logic
	if len(ec.recentBlocks) > MAX_RECENT_BLOCKS {
		// Find and remove the oldest block
		var oldest uint64 = math.MaxUint64
		for blockNum := range ec.recentBlocks {
			if blockNum < oldest {
				oldest = blockNum
			}
		}
		delete(ec.recentBlocks, oldest)
	}
}
