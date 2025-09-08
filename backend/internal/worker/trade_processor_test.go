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

func TestCalculateMetToYouPrice(t *testing.T) {
	tests := []struct {
		name        string
		tradeEvent  *models.TradeEvent
		expected    string
		expectError bool
	}{
		{
			name: "MET to YOU trade - valid amounts",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "100.0", // 100 MET in
				AmountOut: "150.0", // 150 YOU out
			},
			expected: "1.500000", // 150 YOU / 100 MET = 1.5 YOU per MET
		},
		{
			name: "YOU to MET trade - valid amounts",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "YOU",
				TokenOut:  "MET",
				AmountIn:  "200.0", // 200 YOU in
				AmountOut: "100.0", // 100 MET out
			},
			expected: "2.000000", // 200 YOU / 100 MET = 2.0 YOU per MET
		},
		{
			name: "Small amounts with precision",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "0.001",
				AmountOut: "0.0033",
			},
			expected: "3.300000", // 0.0033 / 0.001 = 3.3
		},
		{
			name: "Large amounts",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "1000000.0",
				AmountOut: "750000.0",
			},
			expected: "0.750000", // 0.75 YOU per MET
		},
		{
			name: "Invalid AmountIn",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "invalid",
				AmountOut: "150.0",
			},
			expectError: true,
		},
		{
			name: "Invalid AmountOut",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "100.0",
				AmountOut: "invalid",
			},
			expectError: true,
		},
		{
			name: "Zero MET amount should fail",
			tradeEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "0",
				AmountOut: "150.0",
			},
			expectError: true, // Division by zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculateMetToYouPrice(tt.tradeEvent)

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
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").Return(nil)
				db.On("CreateTrade", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "cache failure",
			inputEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "100.0",
				AmountOut: "150.0",
			},
			setupMocks: func(cache *storageMock.PoolCache, db *storageMock.DB) {
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").
					Return(errors.New("cache connection failed"))
			},
			expectError:   true,
			expectedError: "failed to update cache with trade price",
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
				cache.On("SetPrice", mock.Anything, MET_YOU_PAIR, "1.500000").Return(nil)
				db.On("CreateTrade", mock.Anything).Return(errors.New("db connection failed"))
			},
			expectError:   true,
			expectedError: "failed to store trade event in database",
		},
		{
			name: "invalid price calculation",
			inputEvent: &models.TradeEvent{
				TokenIn:   "MET",
				TokenOut:  "YOU",
				AmountIn:  "invalid",
				AmountOut: "150.0",
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
			err = worker.handleTradeEvent([]byte("test-key"), eventBytes)

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

	err := worker.handleTradeEvent([]byte("key"), invalidJSON)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal trade history event")
}

// Benchmark tests for price calculation
func BenchmarkCalculateMetToYouPrice(b *testing.B) {
	tradeEvent := &models.TradeEvent{
		TokenIn:   "MET",
		TokenOut:  "YOU",
		AmountIn:  "1000.123456789",
		AmountOut: "1500.987654321",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculateMetToYouPrice(tradeEvent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Property-based test for price calculation consistency
func TestPriceCalculationProperties(t *testing.T) {
	// Test that MET->YOU and YOU->MET trades with same rates produce consistent prices
	metAmount := decimal.NewFromFloat(100.0)
	youAmount := decimal.NewFromFloat(150.0)

	// MET to YOU trade
	tradeEvent1 := &models.TradeEvent{
		TokenIn:   "MET",
		TokenOut:  "YOU",
		AmountIn:  metAmount.String(),
		AmountOut: youAmount.String(),
	}

	// YOU to MET trade (inverse)
	tradeEvent2 := &models.TradeEvent{
		TokenIn:   "YOU",
		TokenOut:  "MET",
		AmountIn:  youAmount.String(),
		AmountOut: metAmount.String(),
	}

	price1, err1 := calculateMetToYouPrice(tradeEvent1)
	price2, err2 := calculateMetToYouPrice(tradeEvent2)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Both should give same MET:YOU price
	assert.Equal(t, price1, price2)

	// Verify the actual calculation
	expectedPrice := youAmount.Div(metAmount).StringFixed(6)
	assert.Equal(t, expectedPrice, price1)
}
