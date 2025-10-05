# Token Swap DApp

A simple A simple AMM (Automated Market Maker) token swap application built to learn about smart contracts and blockchain event processing. Supports swapping between two ERC-20 tokens (MET and YOU) using a constant product formula.

## What I Built

This project implements a full-stack decentralized exchange with:

- **Smart Contracts**: Solidity contracts for a liquidity pool and token swapping
- **Event Processing Pipeline**: Go services that monitor blockchain events and process them asynchronously
- **Real-time Data**: Redis caching for current prices and pool reserves
- **Historical Analytics**: PostgreSQL storage for trade history and analytics
- **REST API**: Backend API for frontend integration (work in progress)

## Architecture

The system uses an event-driven architecture:

1. Smart contract emits `Swap` and `Sync` events on-chain
2. Event Listener service monitors the blockchain via WebSocket
3. Events are published to Kafka topics
4. Worker services consume events and update the database and cache
5. API server provides access to real-time and historical data

## Tech Stack

- **Backend**: Go 1.24
- **Smart Contracts**: Solidity + Foundry
- **Message Queue**: Kafka
- **Cache**: Redis
- **Database**: PostgreSQL
- **Blockchain**: Ethereum (tested on Anvil local testnet)

## Technical Decisions

### Kafka for Event Processing
I used Kafka to decouple event detection from processing. This lets me scale workers horizontally and provides built-in message persistence. Events are published to separate topics (`trade-history` and `pool-stats`) for different types of updates.

### Redis for Caching
Current prices and reserves are cached in Redis with 5-minute TTLs. This gives sub-millisecond response times for the most frequently accessed data without constantly querying the database.

### Blockchain Reorg Handling
The system detects chain reorganizations by comparing block hashes. When a conflict is detected (same block number, different hash), it queries the current chain state to determine which events are canonical and rolls back invalid data.

> **Note**: The reorg detection has a known limitation with out-of-order event delivery - it may incorrectly delete valid events if higher block numbers arrive before lower ones during a reorg. A production system would need more robust canonical state tracking.

## Running Locally

### Prerequisites
- Go 1.24+
- Docker
- Foundry

### Setup

**1. Start infrastructure:**
```bash
make up  # Postgres, Kafka, Redis
anvil    # Local Ethereum testnet
```

**2. Run database migrations:**
```bash
go run main.go migrate --config ./config-server.yaml
```

**3. Start services:**
```bash
go run main.go event-listener --config ./config-events.yaml
go run main.go worker --config ./config-worker.yaml
go run main.go server --config ./config-server.yaml
```

**4. Deploy smart contracts:**
```bash
export RPC_URL=http://localhost:8545
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cd contracts
forge script script/Deploy.s.sol --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
```

**5. Execute test swaps:**
```bash
forge script script/Swaps.s.sol --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
```

### Testing

**Unit tests:**
```bash
go test ./...
```

**Integration tests:**
```bash
make test-integration-full
```
