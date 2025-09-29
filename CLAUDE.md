# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a token swap dApp with a Go backend consisting of multiple services that process blockchain events and provide API endpoints. The system uses Kafka for event streaming and Redis for caching.

## Architecture

The backend follows a multi-service architecture:
- **Event Listener**: Monitors blockchain events and publishes to Kafka
- **Worker**: Consumes Kafka events and processes trade/reserve data
- **API Server**: REST API for frontend interactions
- **Migration Tool**: Database schema management

### Service Components

- `cmd/`: CLI commands for each service (server, worker, event_listener, migrate)
- `internal/config/`: Configuration management for different services
- `internal/worker/`: Event processing logic (trades, reserves)  
- `internal/server/`: HTTP handlers and API logic
- `internal/storage/`: Data access layer (Postgres, Redis)
- `internal/kafka/`: Kafka producer/consumer abstractions
- `internal/events/`: Event definitions and handling

## Development Commands

### Infrastructure Management
```bash
# Start all services (Postgres, Kafka, Redis)
make up

# Stop all services  
make down

# View service logs
make logs

# Clean up volumes and data
make clean
```

### Backend Services
```bash
# Build and run API server
go run main.go server --config config-server.yaml

# Build and run event worker
go run main.go worker --config config-worker.yaml  

# Run database migrations
go run main.go migrate --config config-server.yaml

# Build for production
go build -o token-swap main.go
```

### Database Operations
```bash
# Run migrations
go run main.go migrate --config config-server.yaml
```

### Kafka Management
```bash
# Create topics
make create-topics

# List topics  
make list-topics

# Consume trade messages
make consume-trades
```

## Configuration

Each service uses YAML configuration files:
- `config-server.yaml`: API server settings
- `config-worker.yaml`: Worker service settings  
- `config-events.yaml`: Event listener settings

Configuration supports environment variable overrides using Viper.

## Key Dependencies

- **Gin**: HTTP web framework
- **go-ethereum**: Ethereum client libraries
- **confluent-kafka-go**: Kafka client
- **pgx**: PostgreSQL driver
- **go-redis**: Redis client
- **cobra**: CLI framework
- **viper**: Configuration management
- **zerolog**: Structured logging

## Database

PostgreSQL with migrations in `internal/storage/postgres/migrations/sql/`. The system tracks:
- Trade history and analytics
- Pool reserves and liquidity data

## Testing

No specific test framework is configured. Check for test files using standard Go testing patterns.