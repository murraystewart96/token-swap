package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/murraystewart96/token-swap/internal/models"
)

const (
	periodFormat = "2006-01-02 15:04"
)

func (db *DB) GetVolumeAnalytics(start, end time.Time, token string) (*models.VolumeResponse, error) {
	var query string
	var args []any

	if token == "all" {
		query = `
        SELECT 
            COALESCE(SUM(CASE WHEN token_in = 'MET' THEN CAST(amount_in AS NUMERIC) 
                         WHEN token_out = 'MET' THEN CAST(amount_out AS NUMERIC) 
                         ELSE 0 END), 0) as met_volume,
            COALESCE(SUM(CASE WHEN token_in = 'YOU' THEN CAST(amount_in AS NUMERIC) 
                         WHEN token_out = 'YOU' THEN CAST(amount_out AS NUMERIC) 
                         ELSE 0 END), 0) as you_volume,
            COUNT(*) as trade_count
        FROM trades 
        WHERE to_timestamp(timestamp) BETWEEN $1 AND $2`
		args = []any{start, end}
	} else {
		query = `
        SELECT 
            COALESCE(SUM(CASE WHEN token_in = $1 THEN CAST(amount_in AS NUMERIC)
                         WHEN token_out = $1 THEN CAST(amount_out AS NUMERIC)
                         ELSE 0 END), 0) as volume,
            COUNT(*) as trade_count
        FROM trades 
        WHERE (token_in = $1 OR token_out = $1) AND to_timestamp(timestamp) BETWEEN $2 AND $3`
		args = []any{token, start, end}
	}

	row := db.pool.QueryRow(context.Background(), query, args...)

	response := &models.VolumeResponse{Period: fmt.Sprintf("%v to %v", start.Format(periodFormat), end.Format(periodFormat))}

	if token == "all" {
		var metVolume, youVolume string
		err := row.Scan(&metVolume, &youVolume, &response.TradeCount)
		if err != nil {
			return nil, err
		}
		response.TotalVolume.MET = metVolume
		response.TotalVolume.YOU = youVolume
	} else {
		var volume string
		err := row.Scan(&volume, &response.TradeCount)
		if err != nil {
			return nil, err
		}
		if token == "MET" {
			response.TotalVolume.MET = volume
			response.TotalVolume.YOU = "0"
		} else {
			response.TotalVolume.YOU = volume
			response.TotalVolume.MET = "0"
		}
	}

	return response, nil
}

func (db *DB) GetPriceHistory(start, end time.Time, interval time.Duration) (*models.PriceHistoryResponse, error) {
	// Create time buckets and get the last trade in each bucket
	intervalSeconds := int(interval.Seconds())

	query := `
        WITH time_buckets AS (
            SELECT 
                (EXTRACT(EPOCH FROM to_timestamp(timestamp))::bigint / $3) * $3 as bucket_start,
                token_in,
                amount_in,
                amount_out,
                ROW_NUMBER() OVER (
                    PARTITION BY (EXTRACT(EPOCH FROM to_timestamp(timestamp))::bigint / $3) 
                    ORDER BY timestamp DESC
                ) as rn
            FROM trades
            WHERE to_timestamp(timestamp) BETWEEN $1 AND $2
        ),
        price_points AS (
            SELECT 
                bucket_start,
                CASE 
                    WHEN token_in = 'MET' THEN CAST(amount_out AS NUMERIC) / CAST(amount_in AS NUMERIC)
                    ELSE CAST(amount_in AS NUMERIC) / CAST(amount_out AS NUMERIC)
                END as price,
                COALESCE(CAST(amount_in AS NUMERIC), 0) + COALESCE(CAST(amount_out AS NUMERIC), 0) as volume
            FROM time_buckets 
            WHERE rn = 1
        )
        SELECT bucket_start, price, SUM(volume) as total_volume
        FROM price_points
        GROUP BY bucket_start, price
        ORDER BY bucket_start ASC`

	rows, err := db.pool.Query(context.Background(), query, start, end, intervalSeconds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []models.PricePoint
	for rows.Next() {
		var bucketStart int64
		var price, volume string

		err := rows.Scan(&bucketStart, &price, &volume)
		if err != nil {
			return nil, err
		}

		dataPoints = append(dataPoints, models.PricePoint{
			Timestamp: bucketStart,
			Price:     price,
			Volume:    volume,
		})
	}

	return &models.PriceHistoryResponse{
		Period:     fmt.Sprintf("%v to %v", start.Format(periodFormat), end.Format(periodFormat)),
		Interval:   interval.String(),
		DataPoints: dataPoints,
	}, rows.Err()
}

func (db *DB) GetActivityAnalytics(start, end time.Time) (*models.ActivityResponse, error) {
	// Get basic stats
	basicQuery := `
        SELECT 
            COUNT(*) as total_trades,
            COUNT(DISTINCT sender) as unique_traders
        FROM trades 
        WHERE to_timestamp(timestamp) BETWEEN $1 AND $2`

	var totalTrades, uniqueTraders int64
	err := db.pool.QueryRow(context.Background(), basicQuery, start, end).Scan(&totalTrades, &uniqueTraders)
	if err != nil {
		return nil, err
	}

	// Get hourly distribution
	hourlyQuery := `
        SELECT 
            EXTRACT(HOUR FROM to_timestamp(timestamp)) as hour,
            COUNT(*) as trades_count
        FROM trades 
        WHERE to_timestamp(timestamp) BETWEEN $1 AND $2
        GROUP BY EXTRACT(HOUR FROM to_timestamp(timestamp))
        ORDER BY trades_count DESC
        LIMIT 1`

	var peakHour int
	var peakTrades int64
	err = db.pool.QueryRow(context.Background(), hourlyQuery, start, end).Scan(&peakHour, &peakTrades)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	// Calculate average per hour
	duration := end.Sub(start)
	hours := duration.Hours()
	var averagePerHour float64
	if hours > 0 {
		averagePerHour = float64(totalTrades) / hours
	}

	return &models.ActivityResponse{
		Period:         fmt.Sprintf("%v to %v", start.Format(periodFormat), end.Format(periodFormat)),
		TotalTrades:    totalTrades,
		UniqueTraders:  uniqueTraders,
		AveragePerHour: averagePerHour,
		PeakHour: struct {
			Hour   int   `json:"hour"`
			Trades int64 `json:"trades"`
		}{
			Hour:   peakHour,
			Trades: peakTrades,
		},
	}, nil
}
