package server

// type ReservesResponse struct {
// 	METAmount string `json:"met_amount"`
// 	YOUAmount string `json:"you_amount"`
// }

// type CurrentPriceResponse struct {
// 	Price string `json:"current_price"`
// }

// type TradesResponse struct {
// 	Trades     []*models.TradeEvent `json:"trades"`
// 	NextCursor *string              `json:"next_cursor,omitempty"`
// 	HasMore    bool                 `json:"has_more"`
// }

// type VolumeResponse struct {
// 	Period      string `json:"period"`
// 	TotalVolume struct {
// 		MET string `json:"met"`
// 		YOU string `json:"you"`
// 	} `json:"total_volume"`
// 	TradeCount int64 `json:"trade_count"`
// }

// type PriceHistoryResponse struct {
// 	Period     string       `json:"period"`
// 	Interval   string       `json:"interval"`
// 	DataPoints []PricePoint `json:"data_points"`
// }

// type PricePoint struct {
// 	Timestamp int64  `json:"timestamp"`
// 	Price     string `json:"price"`
// 	Volume    string `json:"volume"` // Volume at this time point
// }

// type ActivityResponse struct {
// 	Period         string  `json:"period"`
// 	TotalTrades    int64   `json:"total_trades"`
// 	UniqueTraders  int64   `json:"unique_traders"`
// 	AveragePerHour float64 `json:"average_per_hour"`
// 	PeakHour       struct {
// 		Hour   int   `json:"hour"`
// 		Trades int64 `json:"trades"`
// 	} `json:"peak_hour"`
// }
