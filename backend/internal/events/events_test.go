package events

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/contracts"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDetermineTradeDirection(t *testing.T) {
	tests := []struct {
		name              string
		swapEvent         *contracts.PoolSwap
		expectedTokenIn   string
		expectedTokenOut  string
		expectedAmountIn  string
		expectedAmountOut string
	}{
		{
			name: "MET to YOU trade",
			swapEvent: &contracts.PoolSwap{
				MeTokenIn:   big.NewInt(1000000), // 1 MET (assuming 6 decimals)
				YouTokenOut: big.NewInt(1500000), // 1.5 YOU (assuming 6 decimals)
				YouTokenIn:  big.NewInt(0),
				MeTokenOut:  big.NewInt(0),
			},
			expectedTokenIn:   "MET",
			expectedTokenOut:  "YOU",
			expectedAmountIn:  "1000000",
			expectedAmountOut: "1500000",
		},
		{
			name: "YOU to MET trade",
			swapEvent: &contracts.PoolSwap{
				MeTokenIn:   big.NewInt(0),
				YouTokenOut: big.NewInt(0),
				YouTokenIn:  big.NewInt(2000000), // 2 YOU
				MeTokenOut:  big.NewInt(1000000), // 1 MET
			},
			expectedTokenIn:   "YOU",
			expectedTokenOut:  "MET",
			expectedAmountIn:  "2000000",
			expectedAmountOut: "1000000",
		},
		{
			name: "zero amounts - should default to YOU to MET",
			swapEvent: &contracts.PoolSwap{
				MeTokenIn:   big.NewInt(0),
				YouTokenOut: big.NewInt(0),
				YouTokenIn:  big.NewInt(0),
				MeTokenOut:  big.NewInt(0),
			},
			expectedTokenIn:   "YOU",
			expectedTokenOut:  "MET",
			expectedAmountIn:  "0",
			expectedAmountOut: "0",
		},
		{
			name: "large amounts",
			swapEvent: &contracts.PoolSwap{
				MeTokenIn:   big.NewInt(0).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}),
				YouTokenOut: big.NewInt(0).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}),
				YouTokenIn:  big.NewInt(0),
				MeTokenOut:  big.NewInt(0),
			},
			expectedTokenIn:   "MET",
			expectedTokenOut:  "YOU",
			expectedAmountIn:  big.NewInt(0).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}).String(),
			expectedAmountOut: big.NewInt(0).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}).String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenIn, tokenOut, amountIn, amountOut := determineTradeDirection(tt.swapEvent)

			assert.Equal(t, tt.expectedTokenIn, tokenIn)
			assert.Equal(t, tt.expectedTokenOut, tokenOut)
			assert.Equal(t, tt.expectedAmountIn, amountIn)
			assert.Equal(t, tt.expectedAmountOut, amountOut)
		})
	}
}

// Mock interfaces for testing event handlers
type MockProducer struct {
	mock.Mock
	messages []ProducedMessage
}

type ProducedMessage struct {
	topic string
	key   []byte
	value []byte
}

func (m *MockProducer) Produce(topic string, key, value []byte) error {
	args := m.Called(topic, key, value)
	if args.Error(0) == nil {
		m.messages = append(m.messages, ProducedMessage{
			topic: topic,
			key:   key,
			value: value,
		})
	}
	return args.Error(0)
}

func (m *MockProducer) GetMessages() []ProducedMessage {
	return m.messages
}

func (m *MockProducer) Reset() {
	m.messages = nil
}

func (m *MockProducer) Close() error {
	return nil
}

// Mock pool contract for testing
type MockPoolContract struct {
	mock.Mock
}

func (m *MockPoolContract) ParseSwap(log types.Log) (*contracts.PoolSwap, error) {
	args := m.Called(log)
	return args.Get(0).(*contracts.PoolSwap), args.Error(1)
}

func (m *MockPoolContract) ParseSync(log types.Log) (*contracts.PoolSync, error) {
	args := m.Called(log)
	return args.Get(0).(*contracts.PoolSync), args.Error(1)
}

func (m *MockPoolContract) GetReserves(opts *bind.CallOpts) (struct {
	MeTokenReserve  *big.Int
	YouTokenReserve *big.Int
}, error) {
	args := m.Called(opts)
	return args.Get(0).(struct {
		MeTokenReserve  *big.Int
		YouTokenReserve *big.Int
	}), args.Error(1)
}

func TestHandleSwapEvent(t *testing.T) {
	tests := []struct {
		name           string
		swapEvent      *contracts.PoolSwap
		parseError     error
		producerError  error
		expectError    bool
		expectedTopic  string
		expectedFields map[string]interface{}
	}{
		{
			name: "successful swap event handling - MET to YOU",
			swapEvent: &contracts.PoolSwap{
				Sender:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
				To:          common.HexToAddress("0xfedcba0987654321fedcba0987654321fedcba09"),
				MeTokenIn:   big.NewInt(1000000),
				YouTokenOut: big.NewInt(1500000),
				YouTokenIn:  big.NewInt(0),
				MeTokenOut:  big.NewInt(0),
				Raw: types.Log{
					TxHash:         common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abc123"),
					TxIndex:        5,
					BlockNumber:    12345,
					BlockTimestamp: 1234567890,
					Address:        common.HexToAddress("0xpool1234567890abcdef1234567890abcdef12345678"),
				},
			},
			parseError:    nil,
			producerError: nil,
			expectError:   false,
			expectedTopic: config.TradeHistoryTopic,
			expectedFields: map[string]interface{}{
				"tx_hash":           "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abc123",
				"transaction_index": uint(5),
				"block_number":      uint64(12345),
				"timestamp":         int64(1234567890),
				"sender":            "0x1234567890AbcdEF1234567890abCdef12345678",
				"recipient":         "0xfeDCBA0987654321FedCba0987654321fEDcBa09",
				"token_in":          "MET",
				"token_out":         "YOU",
				"amount_in":         "1000000",
				"amount_out":        "1500000",
				"pool_address":      "0xPool1234567890abCdef1234567890abCdef12345678",
				"event_type":        "swap",
			},
		},
		{
			name: "successful swap event handling - YOU to MET",
			swapEvent: &contracts.PoolSwap{
				Sender:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
				To:          common.HexToAddress("0xfedcba0987654321fedcba0987654321fedcba09"),
				MeTokenIn:   big.NewInt(0),
				YouTokenOut: big.NewInt(0),
				YouTokenIn:  big.NewInt(2000000),
				MeTokenOut:  big.NewInt(1000000),
				Raw: types.Log{
					TxHash:         common.HexToHash("0x4567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef456"),
					TxIndex:        3,
					BlockNumber:    12346,
					BlockTimestamp: 1234567891,
					Address:        common.HexToAddress("0xpool2234567890abcdef1234567890abcdef12345678"),
				},
			},
			parseError:    nil,
			producerError: nil,
			expectError:   false,
			expectedTopic: config.TradeHistoryTopic,
			expectedFields: map[string]interface{}{
				"token_in":     "YOU",
				"token_out":    "MET",
				"amount_in":    "2000000",
				"amount_out":   "1000000",
				"event_type":   "swap",
				"tx_hash":      "0x4567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef456",
				"block_number": uint64(12346),
				"pool_address": "0xPool2234567890abCdef1234567890abCdef12345678",
			},
		},
		{
			name:          "parse swap error",
			swapEvent:     nil,
			parseError:    errors.New("failed to parse swap event"),
			producerError: nil,
			expectError:   false, // Function handles parse errors gracefully
		},
		{
			name: "producer error",
			swapEvent: &contracts.PoolSwap{
				Sender:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
				To:          common.HexToAddress("0x4567890abcdef1234567890abcdef1234567890ab"),
				MeTokenIn:   big.NewInt(1000),
				YouTokenOut: big.NewInt(1500),
				YouTokenIn:  big.NewInt(0),
				MeTokenOut:  big.NewInt(0),
				Raw: types.Log{
					TxHash: common.HexToHash("0x7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456"),
				},
			},
			parseError:    nil,
			producerError: errors.New("kafka producer failed"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProducer := &MockProducer{}
			mockContract := &MockPoolContract{}

			eventLog := &types.Log{
				Topics: []common.Hash{SwapEventSignature},
			}

			// Setup expectations
			if tt.parseError != nil {
				mockContract.On("ParseSwap", *eventLog).Return((*contracts.PoolSwap)(nil), tt.parseError)
			} else if tt.swapEvent != nil {
				mockContract.On("ParseSwap", *eventLog).Return(tt.swapEvent, nil)
				mockProducer.On("Produce", config.TradeHistoryTopic, []byte(tt.swapEvent.Raw.TxHash.Hex()), mock.Anything).Return(tt.producerError)
			}

			// Create event client with mocks
			ec := &EventClient{
				producer:     mockProducer,
				poolContract: mockContract,
			}

			// Execute
			err := ec.handleSwapEvent(eventLog)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify producer was called with correct data
			if tt.expectedTopic != "" && tt.producerError == nil {
				assert.Len(t, mockProducer.GetMessages(), 1)
				message := mockProducer.GetMessages()[0]
				assert.Equal(t, tt.expectedTopic, message.topic)

				// Verify JSON structure
				var tradeEvent models.TradeEvent
				err := json.Unmarshal(message.value, &tradeEvent)
				require.NoError(t, err)

				// Verify expected fields
				for field, expected := range tt.expectedFields {
					switch field {
					case "tx_hash":
						assert.Equal(t, expected, tradeEvent.TxHash)
					case "token_in":
						assert.Equal(t, expected, tradeEvent.TokenIn)
					case "token_out":
						assert.Equal(t, expected, tradeEvent.TokenOut)
					case "amount_in":
						assert.Equal(t, expected, tradeEvent.AmountIn)
					case "amount_out":
						assert.Equal(t, expected, tradeEvent.AmountOut)
					case "event_type":
						assert.Equal(t, expected, tradeEvent.EventType)
					case "block_number":
						assert.Equal(t, expected, tradeEvent.BlockNumber)
					}
				}
			}

			mockContract.AssertExpectations(t)
			mockProducer.AssertExpectations(t)
		})
	}
}

func TestHandleSyncEvent(t *testing.T) {
	tests := []struct {
		name           string
		syncEvent      *contracts.PoolSync
		parseError     error
		producerError  error
		expectError    bool
		expectedTopic  string
		expectedFields map[string]interface{}
	}{
		{
			name: "successful sync event handling",
			syncEvent: &contracts.PoolSync{
				MeTokenAmount:  big.NewInt(5000000),
				YouTokenAmount: big.NewInt(7500000),
				Raw: types.Log{
					TxHash:      common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abc123"),
					BlockNumber: 12347,
					Address:     common.HexToAddress("0xpoolsync567890abcdef1234567890abcdef12345678"),
				},
			},
			parseError:    nil,
			producerError: nil,
			expectError:   false,
			expectedTopic: config.ReserveHistoryTopic,
			expectedFields: map[string]interface{}{
				"tx_hash":      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abc123",
				"block_number": uint64(12347),
				"met_reserve":  "5000000",
				"you_reserve":  "7500000",
				"pool_address": "0xPoolsync567890abCdef1234567890abCdef12345678",
				"event_type":   "sync",
			},
		},
		{
			name:          "parse sync error",
			syncEvent:     nil,
			parseError:    errors.New("failed to parse sync event"),
			producerError: nil,
			expectError:   false, // Function handles parse errors gracefully
		},
		{
			name: "producer error",
			syncEvent: &contracts.PoolSync{
				MeTokenAmount:  big.NewInt(1000),
				YouTokenAmount: big.NewInt(1500),
				Raw: types.Log{
					TxHash: common.HexToHash("0x890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567"),
				},
			},
			parseError:    nil,
			producerError: errors.New("kafka producer failed"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProducer := &MockProducer{}
			mockContract := &MockPoolContract{}

			eventLog := &types.Log{
				Topics: []common.Hash{SyncEventSignature},
			}

			// Setup expectations
			if tt.parseError != nil {
				mockContract.On("ParseSync", *eventLog).Return((*contracts.PoolSync)(nil), tt.parseError)
			} else if tt.syncEvent != nil {
				mockContract.On("ParseSync", *eventLog).Return(tt.syncEvent, nil)
				mockProducer.On("Produce", config.ReserveHistoryTopic, []byte(tt.syncEvent.Raw.TxHash.Hex()), mock.Anything).Return(tt.producerError)
			}

			// Create event client with mocks
			ec := &EventClient{
				producer:     mockProducer,
				poolContract: mockContract,
			}

			// Execute
			err := ec.handleSyncEvent(eventLog)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify producer was called with correct data
			if tt.expectedTopic != "" && tt.producerError == nil {
				assert.Len(t, mockProducer.GetMessages(), 1)
				message := mockProducer.GetMessages()[0]
				assert.Equal(t, tt.expectedTopic, message.topic)

				// Verify JSON structure
				var reserveEvent models.ReserveEvent
				err := json.Unmarshal(message.value, &reserveEvent)
				require.NoError(t, err)

				// Verify expected fields
				for field, expected := range tt.expectedFields {
					switch field {
					case "tx_hash":
						assert.Equal(t, expected, reserveEvent.TxHash)
					case "met_reserve":
						assert.Equal(t, expected, reserveEvent.METReserve)
					case "you_reserve":
						assert.Equal(t, expected, reserveEvent.YOUReserve)
					case "event_type":
						assert.Equal(t, expected, reserveEvent.EventType)
					case "block_number":
						assert.Equal(t, expected, reserveEvent.BlockNumber)
					}
				}
			}

			mockContract.AssertExpectations(t)
			mockProducer.AssertExpectations(t)
		})
	}
}

// Property-based tests
func TestTradeDirectionProperties(t *testing.T) {
	// Test property: exactly one of the input amounts should be > 0
	testCases := []struct {
		name       string
		meTokenIn  *big.Int
		youTokenIn *big.Int
	}{
		{"MET input only", big.NewInt(100), big.NewInt(0)},
		{"YOU input only", big.NewInt(0), big.NewInt(100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			swapEvent := &contracts.PoolSwap{
				MeTokenIn:   tc.meTokenIn,
				YouTokenIn:  tc.youTokenIn,
				YouTokenOut: big.NewInt(50), // Some output
				MeTokenOut:  big.NewInt(50), // Some output
			}

			tokenIn, tokenOut, amountIn, amountOut := determineTradeDirection(swapEvent)

			// Verify one of the directions is chosen
			assert.True(t, (tokenIn == "MET" && tokenOut == "YOU") || (tokenIn == "YOU" && tokenOut == "MET"))

			// Verify amounts are strings of non-negative numbers
			assert.Regexp(t, `^\d+$`, amountIn)
			assert.Regexp(t, `^\d+$`, amountOut)
		})
	}
}
