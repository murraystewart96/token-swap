package worker

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/murraystewart96/token-swap/internal/models"
	storageMock "github.com/murraystewart96/token-swap/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleReserveEvent(t *testing.T) {
	tests := []struct {
		name          string
		inputEvent    *models.ReserveEvent
		setupMocks    func(*storageMock.PoolCache, *storageMock.DB)
		expectError   bool
		expectedError string
	}{
		{
			name: "successful reserve event processing",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				BlockNumber: 12345,
				METReserve:  "1000000.0",
				YOUReserve:  "1500000.0",
				PoolAddress: "0xpool123",
				EventType:   "sync",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				expectedReserves := &models.PoolReserves{
					METAmount: "1000000.0",
					YOUAmount: "1500000.0",
				}
				cache.On("SetReserves", mock.Anything, "0xpool123", expectedReserves).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "cache failure",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "1000000.0",
				YOUReserve:  "1500000.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetReserves", mock.Anything, "0xpool123", mock.Anything).
					Return(errors.New("redis connection failed"))
			},
			expectError:   true,
			expectedError: "failed to update reserves cache",
		},
		{
			name: "database failure",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "1000000.0",
				YOUReserve:  "1500000.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetReserves", mock.Anything, "0xpool123", mock.Anything).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(errors.New("postgres connection failed"))
			},
			expectError:   true,
			expectedError: "failed to store reserve event in database",
		},
		{
			name: "zero reserves",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "0",
				YOUReserve:  "0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				expectedReserves := &models.PoolReserves{
					METAmount: "0",
					YOUAmount: "0",
				}
				cache.On("SetReserves", mock.Anything, "0xpool123", expectedReserves).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "very large reserves",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "999999999999999999999999",
				YOUReserve:  "888888888888888888888888",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				expectedReserves := &models.PoolReserves{
					METAmount: "999999999999999999999999",
					YOUAmount: "888888888888888888888888",
				}
				cache.On("SetReserves", mock.Anything, "0xpool123", expectedReserves).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockCache := &storageMock.PoolCache{}
			DB := &storageMock.DB{}
			tt.setupMocks(mockCache, DB)

			// Create worker with mocks
			worker := &Worker{
				poolCache: mockCache,
				db:        DB,
			}

			// Marshal event to JSON (simulating Kafka message)
			eventBytes, err := json.Marshal(tt.inputEvent)
			require.NoError(t, err)

			// Execute
			err = worker.handleReserveEvent([]byte("test-key"), eventBytes)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockCache.AssertExpectations(t)
			DB.AssertExpectations(t)
		})
	}
}

func TestHandleReserveEventInvalidJSON(t *testing.T) {
	mockCache := &storageMock.PoolCache{}
	DB := &storageMock.DB{}

	worker := &Worker{
		poolCache: mockCache,
		db:        DB,
	}

	// Test various invalid JSON scenarios
	testCases := []struct {
		name        string
		invalidJSON []byte
	}{
		{
			name:        "malformed JSON",
			invalidJSON: []byte(`{"met_reserve": 123, "you_reserve":}`),
		},
		{
			name:        "incomplete JSON",
			invalidJSON: []byte(`{"met_reserve": "100"`),
		},
		{
			name:        "non-JSON data",
			invalidJSON: []byte(`this is not json at all`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := worker.handleReserveEvent([]byte("key"), tc.invalidJSON)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to unmarshal reserve event")
		})
	}
}

func TestReserveEventDataMapping(t *testing.T) {
	// Test that ReserveEvent fields are properly mapped to PoolReserves
	inputEvent := &models.ReserveEvent{
		TxHash:      "0xtest123",
		BlockNumber: 98765,
		METReserve:  "500.123",
		YOUReserve:  "750.456",
		PoolAddress: "0xpooltest",
		EventType:   "sync",
	}

	mockCache := &storageMock.PoolCache{}
	DB := &storageMock.DB{}

	// Capture what gets passed to SetReserves
	var capturedReserves *models.PoolReserves
	mockCache.On("SetReserves", mock.Anything, "0xpooltest", mock.Anything).
		Run(func(args mock.Arguments) {
			capturedReserves = args.Get(2).(*models.PoolReserves)
		}).Return(nil)
	DB.On("CreateReserve", mock.Anything).Return(nil)

	worker := &Worker{
		poolCache: mockCache,
		db:        DB,
	}

	eventBytes, _ := json.Marshal(inputEvent)
	err := worker.handleReserveEvent([]byte("key"), eventBytes)

	require.NoError(t, err)

	// Verify the mapping
	assert.Equal(t, inputEvent.METReserve, capturedReserves.METAmount)
	assert.Equal(t, inputEvent.YOUReserve, capturedReserves.YOUAmount)

	mockCache.AssertExpectations(t)
	DB.AssertExpectations(t)
}

// Benchmark test for reserve event processing
func BenchmarkHandleReserveEvent(b *testing.B) {
	mockCache := &storageMock.PoolCache{}
	DB := &storageMock.DB{}

	// Setup mocks to always succeed
	mockCache.On("SetReserves", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	DB.On("CreateReserve", mock.Anything).Return(nil)

	worker := &Worker{
		poolCache: mockCache,
		db:        DB,
	}

	testEvent := &models.ReserveEvent{
		TxHash:      "0xbenchmark",
		BlockNumber: 12345,
		METReserve:  "1000.0",
		YOUReserve:  "1500.0",
		PoolAddress: "0xpool",
		EventType:   "sync",
	}

	eventBytes, _ := json.Marshal(testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := worker.handleReserveEvent([]byte("key"), eventBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
