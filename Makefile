# Variables
COMPOSE_FILE = docker-compose.yml
PROJECT_NAME = token-swap-dapp

# Default target
.DEFAULT_GOAL := help

## Infrastructure Management
.PHONY: up
up: ## Start all services (Postgres, Kafka, Redis)
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) up -d
	@echo "Services starting... waiting for readiness"
	@sleep 10
	@echo "Infrastructure ready:"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Kafka: localhost:9092"
	@echo "  Redis: localhost:6379"

.PHONY: down
down: ## Stop all services
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) down

.PHONY: restart
restart: down up ## Restart all services

.PHONY: logs
logs: ## Show logs for all services
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) logs -f

.PHONY: clean
clean: down ## Stop services and remove all data
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) down -v
	rm -rf ./kafka-data ./redis-data ./postgres-data
	@echo "All services stopped and data removed"

## Database Management
.PHONY: migrate
migrate: ## Run database migrations
	cd backend && go run . migrate --config ./config-server.yaml

.PHONY: db-shell
db-shell: ## Connect to PostgreSQL shell
	docker exec -it postgres psql -U tokenswap -d tokenswap

## Backend Services
.PHONY: server
server: ## Run API server
	cd backend && go run . server --config ./config-server.yaml

.PHONY: worker
worker: ## Run worker service
	cd backend && go run . worker --config ./config-worker.yaml

.PHONY: event-listener
event-listener: ## Run event listener
	cd backend && go run . event-listener --config ./config-events.yaml

## Development Tools
.PHONY: test
test: ## Run unit tests
	cd backend && go test ./...

.PHONY: build
build: ## Build the application
	cd backend && go build -o ../bin/token-swap .

.PHONY: kafka-console
kafka-console: ## Open Kafka console consumer for trade events
	docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic trade-history --from-beginning --property print.key=true --property key.separator=:

.PHONY: redis-cli
redis-cli: ## Connect to Redis CLI
	docker exec -it redis redis-cli

## Integration Testing
.PHONY: test-integration-setup
test-integration-setup: ## Start integration test infrastructure
	@echo "Starting integration test infrastructure..."
	cd backend/tests/integration && docker-compose -f docker-compose.test.yml -p token-swap-integration-tests up -d
	@echo "Waiting for services to be ready..."
	@sleep 3
	@echo "Running database migrations..."
	cd backend && go run . migrate --config ./tests/integration/config-test.yaml
	@echo "Integration test infrastructure ready:"
	@echo "  PostgreSQL: localhost:5433"
	@echo "  Redis: localhost:6380"
	@echo "  Kafka: localhost:9093"
	@echo "  Anvil: localhost:8546"

.PHONY: test-integration-teardown
test-integration-teardown: ## Stop integration test infrastructure
	cd backend/tests/integration && docker-compose -f docker-compose.test.yml -p token-swap-integration-tests down -v

.PHONY: test-integration
test-integration: ## Run integration tests (requires infrastructure)
	cd backend && go test -v ./tests/integration/... -tags=integration -timeout 2m

.PHONY: test-integration-full
test-integration-full: test-integration-setup test-integration test-integration-teardown ## Full integration test cycle

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
