# Token Swap dApp

A decentralized token swapping application implementing a simplified Uniswap-style AMM (Automated Market Maker) for a single trading pair. This project demonstrates real-time event processing, blockchain integration, and scalable backend architecture patterns.

## Overview

Token Swap is a full-stack decentralized application that enables users to swap between two ERC-20 tokens (MET and YOU) using a constant product formula (x * y = k). The system processes blockchain events in real-time, maintains accurate price data, and provides historical analytics.

### Core Features

- **Token Swapping**: AMM-based token swaps using constant product formula
- **Real-time Price Data**: Live price feeds updated from blockchain events
- **Historical Analytics**: Complete trade and liquidity history tracking
- **Blockchain Re-org Detection**: Robust handling of chain reorganizations
- **Scalable Event Processing**: Asynchronous event handling with Kafka
- **High-Performance Caching**: Redis-based caching for real-time data

## Architecture

The system follows a microservices architecture with event-driven communication:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Smart Contract│    │  Event Listener │    │     Kafka       │
│    (MEYOUPool)  │───▶│    Service      │───▶│   Message       │
└─────────────────┘    └─────────────────┘    │    Queue        │
                                              └─────────┬───────┘
┌─────────────────┐    ┌─────────────────┐              │
│   Frontend      │    │   API Server    │              │
│     (React)     │◀──▶│    (REST)       │              │
└─────────────────┘    └─────────────────┘              │
                                ▲                       │
┌─────────────────┐    ┌─────────────────┐              │
│     Redis       │◀──▶│    Worker       │◀─────────────┘
│    Cache        │    │   Services      │
└─────────────────┘    └─────────────────┘
                                ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │   Database      │
                       └─────────────────┘
```

### Components

- **Smart Contracts**: Solidity contracts implementing AMM functionality
- **Event Listener**: Monitors blockchain events and publishes to Kafka
- **Worker Services**: Process events asynchronously and update data stores
- **API Server**: REST API for frontend interactions
- **Redis Cache**: Real-time price and reserve data caching
- **PostgreSQL**: Historical data and analytics storage

## Key Design Decisions

### 1. Kafka Message Queue for Event Processing

**Why Kafka?**
- **Scalability**: Horizontal scaling of event processors through consumer groups
- **Reliability**: Persistent message storage with configurable retention
- **Decoupling**: Loose coupling between event detection and processing
- **Ordering**: Maintains event ordering within partitions for consistency

**Implementation:**
- Separate topics for trade events (`trade-history`) and reserve updates (`pool-stats`)
- Event listener publishes structured JSON messages
- Multiple worker instances can consume events in parallel
- Built-in retry mechanisms and error handling

### 2. Redis Caching for Real-time Data

**Why Redis?**
- **Performance**: Sub-millisecond latency for price queries
- **Memory Efficiency**: Optimized for frequent reads with TTL-based expiration
- **Atomic Operations**: Consistent updates for price and reserve data
- **Scalability**: Easy to scale with Redis Cluster for high throughput

**Caching Strategy:**
- **Price Data**: 5-minute TTL, updated on each Sync event
- **Reserve Data**: 5-minute TTL, cached as JSON objects
- **Cache Keys**: Namespaced by pool address for multi-pool support
- **Fallback**: Database queries when cache misses occur

### 3. Blockchain Re-org Detection

**The Problem:**
Ethereum blockchain can reorganize, invalidating previously confirmed events and requiring data rollback.

**Our Solution:**
- **Conflict Detection**: Compare block hashes for the same block number
- **Canonical Chain Verification**: Query current blockchain state to resolve conflicts
- **Automatic Rollback**: Remove events from reorganized blocks
- **Database Consistency**: Maintain data integrity across chain re-orgs

**Note:** Current implementation has a limitation with out-of-order event delivery during reorgs - may delete valid canonical events if higher block numbers arrive first. A production system would need to track canonical status more robustly.

### 4. Event-Driven Architecture

**Event Flow:**
1. Smart contract emits `Swap` and `Sync` events
2. Event Listener captures events via WebSocket subscription
3. Events published to appropriate Kafka topics with structured payloads
4. Worker services consume events and update databases/caches
5. API server serves real-time and historical data to frontend

## Development Setup

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- Foundry (for smart contract development)

### Quick Start

1. **Start Infrastructure**
   ```bash
   make up  # Starts Postgres, Kafka, and Redis
   anvil    # Run local anvil testnet
   ```

2. **Run Database Migrations**
   ```bash
   go run main.go migrate --config ./config-server.yaml
   ```

3. **Start Backend Services**
   ```bash
   # Terminal 1: Event Listener
   go run main.go event-listener --config ./config-events.yaml

   # Terminal 2: Worker
   go run main.go worker --config ./config-worker.yaml

   # Terminal 3: API Server
   go run main.go server --config ./config-server.yaml
   ```

4. **Deploy Smart Contracts**
   ```bash
    # Set environment variables for local development
   export RPC_URL=http://localhost:8545
   export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

   cd contracts
   forge script script/Deploy.s.sol --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
   ```

5. **Execute Test Swaps** (optional)
   ```bash
   # Execute test swaps
   forge script script/Swaps.s.sol --rpc-url $RPC_URL --broadcast --private-key $PRIVATE_KEY
   ```

### Configuration

Each service uses YAML configuration with environment variable overrides:

- `config-server.yaml`: API server settings
- `config-worker.yaml`: Worker service settings
- `config-events.yaml`: Event listener configuration

### Testing

```bash
# Unit tests
go test ./...

# Integration tests 
make test-integration-full
```





