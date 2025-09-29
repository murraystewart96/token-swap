package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/murraystewart96/token-swap/internal/config"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(cfg *config.DB) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Optional: configure pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) GetConn() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Exec(query string) error {
	_, err := db.pool.Exec(context.Background(), query)
	return err
}

func (db *DB) RollbackEvents(blockNumber uint64) error {
	ctx := context.Background()
	
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tradesQuery := `DELETE FROM trades WHERE block_number >= $1`
	_, err = tx.Exec(ctx, tradesQuery, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to delete trades: %w", err)
	}

	reservesQuery := `DELETE FROM reserves WHERE block_number >= $1`
	_, err = tx.Exec(ctx, reservesQuery, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to delete reserves: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
