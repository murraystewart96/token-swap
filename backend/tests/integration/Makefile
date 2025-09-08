# Integration Test Makefile
COMPOSE_FILE = docker-compose.test.yml
PROJECT_NAME = token-swap-integration-tests

## Setup commands
.PHONY: test-integration-setup
test-integration-setup: ## Start all test infrastructure services
	@echo "Starting test infrastructure..."
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) up -d
	@echo "Waiting for services to be ready..."
	@sleep 10

.PHONY: test-integration-teardown  
test-integration-teardown: ## Stop and remove all test infrastructure
	docker-compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME) down -v

## Test execution
.PHONY: test-integration
test-integration: ## Run all integration tests
	cd .. && go test -v ./integration-tests/... -tags=integration

.PHONY: test-integration-full
test-integration-full: test-integration-setup test-integration test-integration-teardown ## Full test cycle: setup -> test -> teardown