package models

type TradeEvent struct {
	// Transaction identifiers
	TxHash      string `json:"tx_hash"`
	BlockNumber uint64 `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`

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
	EventType   string `json:"event_type"`   // "swap"
}
