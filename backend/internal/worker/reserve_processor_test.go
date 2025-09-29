package worker

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/murraystewart96/token-swap/internal/models"
	storageMock "github.com/murraystewart96/token-swap/internal/storage/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCalculateMetToYouPriceFromReserves(t *testing.T) {
	tests := []struct {
		name        string
		reserves    *models.PoolReserves
		expected    string
		expectError bool
	}{
		{
			name: "Valid reserves - balanced pool",
			reserves: &models.PoolReserves{
				METAmount: "100.0",
				YOUAmount: "150.0",
			},
			expected: "1.500000", // 150 YOU / 100 MET = 1.5 YOU per MET
		},
		{
			name: "Valid reserves - unbalanced pool",
			reserves: &models.PoolReserves{
				METAmount: "200.0",
				YOUAmount: "100.0",
			},
			expected: "0.500000", // 100 YOU / 200 MET = 0.5 YOU per MET
		},
		{
			name: "Small amounts with precision",
			reserves: &models.PoolReserves{
				METAmount: "0.001",
				YOUAmount: "0.0033",
			},
			expected: "3.300000", // 0.0033 / 0.001 = 3.3
		},
		{
			name: "Large amounts",
			reserves: &models.PoolReserves{
				METAmount: "1000000.0",
				YOUAmount: "750000.0",
			},
			expected: "0.750000", // 0.75 YOU per MET
		},
		{
			name: "Invalid MET amount",
			reserves: &models.PoolReserves{
				METAmount: "invalid",
				YOUAmount: "150.0",
			},
			expectError: true,
		},
		{
			name: "Invalid YOU amount",
			reserves: &models.PoolReserves{
				METAmount: "100.0",
				YOUAmount: "invalid",
			},
			expectError: true,
		},
		{
			name: "Zero MET amount should fail",
			reserves: &models.PoolReserves{
				METAmount: "0",
				YOUAmount: "150.0",
			},
			expectError: true, // Division by zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculateMetToYouPrice(tt.reserves)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

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
				METReserve:  "100.0",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").Return(nil)
				expectedReserves := &models.PoolReserves{
					METAmount: "100.0",
					YOUAmount: "150.0",
				}
				cache.On("SetReserves", mock.Anything, "0xpool123", expectedReserves).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "price cache failure",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "100.0",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").
					Return(errors.New("cache connection failed"))
			},
			expectError:   true,
			expectedError: "failed to update cache with trade price",
		},
		{
			name: "reserves cache failure",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "100.0",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").Return(nil)
				cache.On("SetReserves", mock.Anything, "0xpool123", mock.Anything).
					Return(errors.New("cache connection failed"))
			},
			expectError:   true,
			expectedError: "failed to update reserves cache",
		},
		{
			name: "database failure",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "100.0",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").Return(nil)
				cache.On("SetReserves", mock.Anything, "0xpool123", mock.Anything).Return(nil)
				db.On("CreateReserve", mock.Anything).Return(errors.New("postgres connection failed"))
			},
			expectError:   true,
			expectedError: "failed to store reserve event in database",
		},
		{
			name: "invalid price calculation - zero MET reserves",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "0",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks:    func(*storageMock.PoolCache, *storageMock.DB) {}, // No mocks needed, fails before
			expectError:   true,
			expectedError: "failed to calculate trading price",
		},
		{
			name: "invalid price calculation - invalid MET reserves",
			inputEvent: &models.ReserveEvent{
				TxHash:      "0xabc123",
				METReserve:  "invalid",
				YOUReserve:  "150.0",
				PoolAddress: "0xpool123",
			},
			setupMocks:    func(*storageMock.PoolCache, *storageMock.DB) {}, // No mocks needed, fails before
			expectError:   true,
			expectedError: "failed to calculate trading price",
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
			err = worker.handleReserveEvent(t.Context(), []byte("test-key"), eventBytes)

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
			err := worker.handleReserveEvent(t.Context(), []byte("key"), tc.invalidJSON)

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
	err := worker.handleReserveEvent(t.Context(), []byte("key"), eventBytes)

	require.NoError(t, err)

	// Verify the mapping
	assert.Equal(t, inputEvent.METReserve, capturedReserves.METAmount)
	assert.Equal(t, inputEvent.YOUReserve, capturedReserves.YOUAmount)

	mockCache.AssertExpectations(t)
	DB.AssertExpectations(t)
}

// Benchmark tests for price calculation from reserves
func BenchmarkCalculateMetToYouPriceFromReserves(b *testing.B) {
	reserves := &models.PoolReserves{
		METAmount: "1000.123456789",
		YOUAmount: "1500.987654321",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculateMetToYouPrice(reserves)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Property-based test for price calculation consistency from reserves
func TestPriceCalculationPropertiesFromReserves(b *testing.T) {
	// Test that price calculation is consistent regardless of decimal precision
	metAmount := decimal.NewFromFloat(100.0)
	youAmount := decimal.NewFromFloat(150.0)

	reserves := &models.PoolReserves{
		METAmount: metAmount.String(),
		YOUAmount: youAmount.String(),
	}

	price, err := calculateMetToYouPrice(reserves)
	require.NoError(b, err)

	// Verify the actual calculation
	expectedPrice := youAmount.Div(metAmount).StringFixed(6)
	assert.Equal(b, expectedPrice, price)

	// Test with different precision but same ratio
	reserves2 := &models.PoolReserves{
		METAmount: metAmount.Mul(decimal.NewFromInt(2)).String(),
		YOUAmount: youAmount.Mul(decimal.NewFromInt(2)).String(),
	}

	price2, err2 := calculateMetToYouPrice(reserves2)
	require.NoError(b, err2)

	// Should give same price despite different amounts
	assert.Equal(b, price, price2)
}

// Benchmark test for reserve event processing
func BenchmarkHandleReserveEvent(b *testing.B) {
	mockCache := &storageMock.PoolCache{}
	DB := &storageMock.DB{}

	// Setup mocks to always succeed
	mockCache.On("SetPrice", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
	}

	eventBytes, _ := json.Marshal(testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := worker.handleReserveEvent(b.Context(), []byte("key"), eventBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
