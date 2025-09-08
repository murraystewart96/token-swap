-- +goose Up
CREATE TABLE reserves (
    id BIGSERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) NOT NULL,
    transaction_index INTEGER NOT NULL DEFAULT 0,
    block_number BIGINT NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    met_reserve NUMERIC(78, 0) NOT NULL,
    you_reserve NUMERIC(78, 0) NOT NULL,
    pool_address VARCHAR(42) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS reserves;

