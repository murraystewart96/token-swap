package events

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	storageMock "github.com/murraystewart96/token-swap/internal/storage/mock"
	ethMock "github.com/murraystewart96/token-swap/pkg/eth/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to create test event
func createTestEvent(blockNumber uint64, blockHash common.Hash) *types.Log {
	return &types.Log{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
		Topics:      []common.Hash{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")}, // dummy topic
	}
}

// Helper to create test block - we'll create different blocks that naturally have different hashes
func createTestBlock(blockNumber uint64, blockHash common.Hash) *types.Block {
	// Create headers with different content to get different hashes
	header := &types.Header{
		Number:     big.NewInt(int64(blockNumber)),
		ParentHash: blockHash,             // Use the desired hash as parent hash to make blocks different
		GasLimit:   uint64(blockHash[31]), // Use last byte of hash as gas limit for uniqueness
	}
	return types.NewBlockWithHeader(header)
}

func TestCheckForChainReorg_NoConflict_FirstTimeSeeingBlock(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	ec := &EventClient{
		ethClient:    mockClient,
		db:           mockDB,
		recentBlocks: make(map[uint64]common.Hash),
	}

	blockHash := common.HexToHash("0xa1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456")
	event := createTestEvent(100, blockHash)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.NoError(t, err)
	assert.False(t, reorg)
	assert.Equal(t, blockHash, ec.recentBlocks[100])

	// Verify no external calls were made
	mockClient.AssertNotCalled(t, "BlockByNumber")
	mockDB.AssertNotCalled(t, "RollbackEvents")
}

func TestCheckForChainReorg_NoConflict_SameHashAsStored(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	blockHash := common.HexToHash("0xa1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456")
	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: blockHash, // Already seen this block with same hash
		},
	}

	event := createTestEvent(100, blockHash)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.NoError(t, err)
	assert.False(t, reorg)

	// Verify no external calls were made
	mockClient.AssertNotCalled(t, "BlockByNumber")
	mockDB.AssertNotCalled(t, "RollbackEvents")
}

func TestCheckForChainReorg_DetectsReorg_NewEventIsCanonical(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	oldHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123")
	newHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456")

	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: oldHash, // We've seen block 100 with different hash
		},
	}

	// Mock canonical chain first to get the canonical hash
	canonicalBlock := createTestBlock(100, newHash)
	canonicalHash := canonicalBlock.Hash() // Get the actual computed hash

	// Create event that matches the canonical hash (new event is canonical)
	event := createTestEvent(100, canonicalHash)

	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(100)).Return(canonicalBlock, nil)

	// Mock successful database rollback
	mockDB.On("RollbackEvents", uint64(100)).Return(nil)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.NoError(t, err)
	assert.True(t, reorg)
	assert.Equal(t, canonicalHash, ec.recentBlocks[100]) // Should update to canonical hash

	// Verify external calls
	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCheckForChainReorg_DetectsReorg_ExistingEventIsCanonical(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	// Create canonical block first to get the actual canonical hash
	canonicalBlock := createTestBlock(100, common.Hash{1}) // Some dummy hash for creating different blocks
	canonicalHash := canonicalBlock.Hash()                 // Get the actual computed hash
	nonCanonicalHash := common.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: canonicalHash, // We've seen canonical block with the actual hash
		},
	}

	// New event has different hash (non-canonical)
	event := createTestEvent(100, nonCanonicalHash)

	// Mock canonical chain returning the existing hash (existing is canonical)
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(100)).Return(canonicalBlock, nil)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.NoError(t, err)
	assert.False(t, reorg)                               // No reorg action taken
	assert.Equal(t, canonicalHash, ec.recentBlocks[100]) // Should keep existing hash

	// Verify no database rollback occurred
	mockClient.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "RollbackEvents")
}

func TestCheckForChainReorg_DetectsReorg_BothEventsNonCanonical(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	oldHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123")
	newHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456")
	canonicalHash := common.HexToHash("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")

	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: oldHash,
		},
	}

	event := createTestEvent(100, newHash)

	// Mock canonical chain returning completely different hash
	canonicalBlock := createTestBlock(100, canonicalHash)
	actualCanonicalHash := canonicalBlock.Hash() // Get the actual computed hash
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(100)).Return(canonicalBlock, nil)

	// Mock successful database rollback
	mockDB.On("RollbackEvents", uint64(100)).Return(nil)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.NoError(t, err)
	assert.True(t, reorg)
	assert.Equal(t, actualCanonicalHash, ec.recentBlocks[100]) // Should update to canonical hash

	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCheckForChainReorg_RollbackDatabaseError(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	oldHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123")
	newHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456")

	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: oldHash,
		},
	}

	event := createTestEvent(100, newHash)

	// Mock canonical chain returning new hash
	canonicalBlock := createTestBlock(100, newHash)
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(100)).Return(canonicalBlock, nil)

	// Mock database rollback failure
	mockDB.On("RollbackEvents", uint64(100)).Return(assert.AnError)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.Error(t, err)
	assert.True(t, reorg)

	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCheckForChainReorg_CanonicalChainQueryError(t *testing.T) {
	// Setup
	mockClient := &ethMock.EthClient{}
	mockDB := &storageMock.DB{}

	oldHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123")
	newHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456")

	ec := &EventClient{
		ethClient: mockClient,
		db:        mockDB,
		recentBlocks: map[uint64]common.Hash{
			100: oldHash,
		},
	}

	event := createTestEvent(100, newHash)

	// Mock canonical chain query failure
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(100)).Return(nil, assert.AnError)

	// Execute
	reorg, err := ec.CheckForChainReorg(t.Context(), event)

	// Assert
	assert.Error(t, err)
	assert.False(t, reorg)
	assert.Contains(t, err.Error(), "failed to get canonical hash")

	mockClient.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "RollbackEvents")
}

// Test the memory management function
func TestAddBlock_MemoryManagement(t *testing.T) {
	ec := &EventClient{
		recentBlocks: make(map[uint64]common.Hash),
	}

	// Fill up to MAX_RECENT_BLOCKS
	for i := uint64(1); i <= MAX_RECENT_BLOCKS; i++ {
		hash := common.HexToHash(fmt.Sprintf("0x%064d", i))
		ec.addBlock(i, hash)
	}

	assert.Equal(t, MAX_RECENT_BLOCKS, len(ec.recentBlocks))

	// Add one more - should trigger cleanup
	newHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000999")
	ec.addBlock(MAX_RECENT_BLOCKS+1, newHash)

	// Should still be at max size
	assert.Equal(t, MAX_RECENT_BLOCKS, len(ec.recentBlocks))

	// Should contain the newest block
	assert.Equal(t, newHash, ec.recentBlocks[MAX_RECENT_BLOCKS+1])

	// Should have removed the oldest block
	_, exists := ec.recentBlocks[1]
	assert.False(t, exists)
}
