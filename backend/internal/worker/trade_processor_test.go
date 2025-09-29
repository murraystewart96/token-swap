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

func TestHandleTradeEvent(t *testing.T) {
	tests := []struct {
		name          string
		inputEvent    *models.TradeEvent
		setupMocks    func(*storageMock.PoolCache, *storageMock.DB)
		expectError   bool
		expectedError string
	}{
		{
			name: "successful trade event processing",
			inputEvent: &models.TradeEvent{
				TxHash:    "0x123",
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "100.0",
				AmountOut: "150.0",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				db.On("CreateTrade", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "database failure",
			inputEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "100.0",
				AmountOut: "150.0",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				db.On("CreateTrade", mock.Anything).Return(errors.New("db connection failed"))
			},
			expectError:   true,
			expectedError: "failed to store trade event in database",
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
			err = worker.handleTradeEvent(t.Context(), []byte("test-key"), eventBytes)

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

func TestHandleTradeEventInvalidJSON(t *testing.T) {
	mockCache := &storageMock.PoolCache{}
	DB := &storageMock.DB{}

	worker := &Worker{
		poolCache: mockCache,
		db:        DB,
	}

	invalidJSON := []byte(`{"invalid": json}`)

	err := worker.handleTradeEvent(t.Context(), []byte("key"), invalidJSON)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal trade history event")
}

// Benchmark test for trade event processing
func BenchmarkHandleTradeEvent(b *testing.B) {
	mockCache := &storageMock.PoolCache{}
	mockDB := &storageMock.DB{}

	// Setup mocks to always succeed
	mockDB.On("CreateTrade", mock.Anything).Return(nil)

	worker := &Worker{
		poolCache: mockCache,
		db:        mockDB,
	}

	testEvent := &models.TradeEvent{
		TxHash:    "0xbenchmark",
		TokenIn:   "MET",
		TokenOut:  "YOU",
		AmountIn:  "1000.0",
		AmountOut: "1500.0",
	}

	eventBytes, _ := json.Marshal(testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := worker.handleTradeEvent(b.Context(), []byte("key"), eventBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
