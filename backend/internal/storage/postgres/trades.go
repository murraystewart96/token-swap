package postgres

import (
	"context"
	"time"

	"github.com/murraystewart96/token-swap/internal/models"
)

func (db DB) CreateTrade(trade *models.TradeEvent) error {
	query := `
        INSERT INTO trades (tx_hash, block_number, transaction_index, timestamp, sender, recipient, 
                           token_in, token_out, amount_in, amount_out, pool_address)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.pool.Exec(context.Background(), query,
		trade.TxHash, trade.BlockNumber, trade.TransactionIndex, trade.Timestamp,
		trade.Sender, trade.Recipient, trade.TokenIn, trade.TokenOut,
		trade.AmountIn, trade.AmountOut, trade.PoolAddress)

	return err
}

func (db *DB) GetTradesByTimeRange(start, end time.Time) ([]*models.TradeEvent, error) {
	query := `
        SELECT tx_hash, block_number, transaction_index, timestamp, sender, recipient,
               token_in, token_out, amount_in, amount_out, pool_address
        FROM trades
        WHERE to_timestamp(timestamp) BETWEEN $1 AND $2
        ORDER BY block_number ASC`

	rows, err := db.pool.Query(context.Background(), query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*models.TradeEvent
	for rows.Next() {
		trade := &models.TradeEvent{}
		err := rows.Scan(&trade.TxHash, &trade.BlockNumber, &trade.TransactionIndex, &trade.Timestamp,
			&trade.Sender, &trade.Recipient, &trade.TokenIn,
			&trade.TokenOut, &trade.AmountIn, &trade.AmountOut,
			&trade.PoolAddress)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, rows.Err()
}

func (db *DB) GetTradesByCursor(cursorBlock uint64, cursorTx uint, limit int) ([]*models.TradeEvent, error) {
	var query string
	var args []any

	if cursorBlock == 0 && cursorTx == 0 {
		// First page - no cursor provided
		query = `
            SELECT tx_hash, block_number, transaction_index, timestamp, sender, recipient,
                   token_in, token_out, amount_in, amount_out, pool_address
            FROM trades
            ORDER BY block_number DESC, transaction_index DESC
            LIMIT $1`
		args = []any{limit}
	} else {
		// Subsequent pages - use cursor
		query = `
            SELECT tx_hash, block_number, transaction_index, timestamp, sender, recipient,
                   token_in, token_out, amount_in, amount_out, pool_address
            FROM trades
            WHERE (block_number < $1) 
               OR (block_number = $1 AND transaction_index < $2)
            ORDER BY block_number DESC, transaction_index DESC
            LIMIT $3`
		args = []any{cursorBlock, cursorTx, limit}
	}

	rows, err := db.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*models.TradeEvent
	for rows.Next() {
		trade := &models.TradeEvent{}
		err := rows.Scan(&trade.TxHash, &trade.BlockNumber, &trade.TransactionIndex,
			&trade.Timestamp, &trade.Sender, &trade.Recipient, &trade.TokenIn,
			&trade.TokenOut, &trade.AmountIn, &trade.AmountOut, &trade.PoolAddress)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, rows.Err()
}

func (db *DB) UpdateConfirmedTrades(confirmationThreshold uint64) error {
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
