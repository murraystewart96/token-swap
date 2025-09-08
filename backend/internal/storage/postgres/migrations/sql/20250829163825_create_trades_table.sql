-- +goose Up
CREATE TABLE trades (
    id BIGSERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) NOT NULL UNIQUE,
    transaction_index INTEGER NOT NULL DEFAULT 0,
    block_number BIGINT NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    timestamp BIGINT NOT NULL,
    sender VARCHAR(42) NOT NULL,
    recipient VARCHAR(42) NOT NULL,
    token_in VARCHAR(10) NOT NULL,
    token_out VARCHAR(10) NOT NULL,
    amount_in NUMERIC(78, 0) NOT NULL,
    amount_out NUMERIC(78, 0) NOT NULL,
    pool_address VARCHAR(42) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS trades;

