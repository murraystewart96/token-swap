package models

type TradeEvent struct {
	// Transaction identifiers
	TxHash           string `json:"tx_hash"`
	BlockNumber      uint64 `json:"block_number"`
	Timestamp        int64  `json:"timestamp"`
	TransactionIndex uint   `json:"transaction_index"`

	// Trade participants
	Sender    string `json:"sender"`    // Who initiated the trade
	Recipient string `json:"recipient"` // Who received the output (usually same as sender)

	// Trade details
	TokenIn   string `json:"token_in"`   // "MET" or "YOU"
	TokenOut  string `json:"token_out"`  // "YOU" or "MET"
	AmountIn  string `json:"amount_in"`  // Input amount
	AmountOut string `json:"amount_out"` // Output amount

	// Context
	PoolAddress string `json:"pool_address"` // Which pool
}

type ReserveEvent struct {
	TxHash      string `json:"tx_hash"`
	BlockNumber uint64 `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`

	METReserve  string `json:"met_reserve"`
	YOUReserve  string `json:"you_reserve"`
	PoolAddress string `json:"pool_address"`
}

type PoolReserves struct {
	METAmount string `json:"met_amount"`
	YOUAmount string `json:"you_amount"`
}

// *** API RESPONSE MODELS ***

type ReservesResponse struct {
	METAmount string `json:"met_amount"`
	YOUAmount string `json:"you_amount"`
}

type CurrentPriceResponse struct {
	Price string `json:"current_price"`
}

type TradesResponse struct {
	Trades     []*TradeEvent `json:"trades"`
	NextCursor *string       `json:"next_cursor,omitempty"`
	HasMore    bool          `json:"has_more"`
}

type VolumeResponse struct {
	Period      string `json:"period"`
	TotalVolume struct {
		MET string `json:"met"`
		YOU string `json:"you"`
	} `json:"total_volume"`
	TradeCount int64 `json:"trade_count"`
}

type PriceHistoryResponse struct {
	Period     string       `json:"period"`
	Interval   string       `json:"interval"`
	DataPoints []PricePoint `json:"data_points"`
}

type PricePoint struct {
	Timestamp int64  `json:"timestamp"`
	Price     string `json:"price"`
	Volume    string `json:"volume"` // Volume at this time point
}

type ActivityResponse struct {
	Period         string  `json:"period"`
	TotalTrades    int64   `json:"total_trades"`
	UniqueTraders  int64   `json:"unique_traders"`
	AveragePerHour float64 `json:"average_per_hour"`
	PeakHour       struct {
		Hour   int   `json:"hour"`
		Trades int64 `json:"trades"`
	} `json:"peak_hour"`
}
