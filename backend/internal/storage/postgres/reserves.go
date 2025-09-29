package postgres

import (
	"context"
	"time"

	"github.com/murraystewart96/token-swap/internal/models"
)

func (db *DB) CreateReserve(reserve *models.ReserveEvent) error {
	query := `
        INSERT INTO reserves (tx_hash, block_number, timestamp, met_reserve, you_reserve, pool_address)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.pool.Exec(context.Background(), query,
		reserve.TxHash,
		reserve.BlockNumber,
		reserve.Timestamp,
		reserve.METReserve,
		reserve.YOUReserve,
		reserve.PoolAddress)

	return err
}

func (db *DB) GetReservesByTimeRange(start, end time.Time) ([]*models.ReserveEvent, error) {
	query := `
        SELECT tx_hash, block_number, timestamp, met_reserve, you_reserve, pool_address
        FROM reserves
        WHERE to_timestamp(timestamp) BETWEEN $1 AND $2
        ORDER BY block_number ASC`

	rows, err := db.pool.Query(context.Background(), query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reserves []*models.ReserveEvent
	for rows.Next() {
		reserve := &models.ReserveEvent{}
		err := rows.Scan(&reserve.TxHash, &reserve.BlockNumber, &reserve.Timestamp,
			&reserve.METReserve, &reserve.YOUReserve,
			&reserve.PoolAddress)
		if err != nil {
			return nil, err
		}
		reserves = append(reserves, reserve)
	}

	return reserves, rows.Err()
}

func (db *DB) UpdateConfirmedReserves(confirmationThreshold uint64) error {
	query := `
        UPDATE reserves 
        SET confirmed = true 
        WHERE confirmed = false 
          AND block_number <= $1`

	_, err := db.pool.Exec(context.Background(), query, confirmationThreshold)
	if err != nil {
		return err
	}

	return nil
}
