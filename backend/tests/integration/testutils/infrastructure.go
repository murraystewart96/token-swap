package testutils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/pkg/kafka"
	"github.com/murraystewart96/token-swap/internal/storage"
	"github.com/murraystewart96/token-swap/internal/storage/postgres"
	"github.com/murraystewart96/token-swap/internal/storage/redis"
)

// TestInfrastructure connects to running test services
// Much simpler - we just connect to services started by docker-compose
type TestInfrastructure struct {
	DB            storage.DB
	PoolCache     storage.PoolCache
	KafkaProducer kafka.IProducer
	KafkaConsumer kafka.IConsumer
	ETHClient     *ethclient.Client
}

// GetEnvWithDefault returns environment variable value or default if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestInfrastructure connects to running test services
func SetupTestInfrastructure(ctx context.Context) (*TestInfrastructure, error) {
	// Connect to PostgreSQL using environment variables
	dbConfig := &config.DB{
		Host:     GetEnvWithDefault("TEST_DB_HOST", "localhost"),
		Port:     GetEnvWithDefault("TEST_DB_PORT", "5433"),
		Name:     GetEnvWithDefault("TEST_DB_NAME", "tokenswap_test"),
		User:     GetEnvWithDefault("TEST_DB_USER", "test_user"),
		Password: GetEnvWithDefault("TEST_DB_PASSWORD", "test_password"),
	}
	db, err := postgres.NewDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Connect to Redis
	redisConfig := &config.Redis{
		Addr: GetEnvWithDefault("TEST_REDIS_ADDR", "localhost:6380"),
	}
	poolCache := redis.NewCache(redisConfig)

	// Connect to Kafka producer
	producerConfig := &config.KafkaProducer{
		BootstrapServers: GetEnvWithDefault("TEST_KAFKA_BROKER", "localhost:9093"),
		Acks:             GetEnvWithDefault("TEST_KAFKA_ACKS", "all"),
	}
	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// Connect to Kafka consumer
	consumerConfig := &config.KafkaConsumer{
		BootstrapServers: GetEnvWithDefault("TEST_KAFKA_BROKER", "localhost:9093"),
		GroupID:          fmt.Sprintf("%s-%d", GetEnvWithDefault("TEST_KAFKA_GROUP_ID", "test-consumer-group"), time.Now().UnixNano()),
		OffsetReset:      GetEnvWithDefault("TEST_KAFKA_OFFSET_RESET", "earliest"),
	}
	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// Connect to Ethereum client
	ethNodeURL := GetEnvWithDefault("TEST_ETH_NODE_URL", "http://localhost:8546")
	ethClient, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum client: %w", err)
	}

	infra := &TestInfrastructure{
		DB:            db,
		PoolCache:     poolCache,
		KafkaProducer: producer,
		KafkaConsumer: consumer,
		ETHClient:     ethClient,
	}

	return infra, nil
}

// Reset clears all data from services, useful between tests
func (ti *TestInfrastructure) Reset() error {
	// Reset Redis cache
	if err := ti.PoolCache.Reset(); err != nil {
		return fmt.Errorf("failed to reset redis cache: %w", err)
	}

	// Reset database - truncate all tables
	if err := ti.resetDatabase(); err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	return nil
}

// resetDatabase truncates all tables to clean state
func (ti *TestInfrastructure) resetDatabase() error {
	query := `
		TRUNCATE TABLE trades RESTART IDENTITY CASCADE;
		TRUNCATE TABLE reserves RESTART IDENTITY CASCADE;
	`

	// Execute the truncate query
	// Note: This is simplified - you'll need to adapt based on your actual DB interface
	return ti.DB.Exec(query)
}

// Cleanup closes all service connections
func (ti *TestInfrastructure) Cleanup() error {
	// Close database connection
	ti.DB.Close()

	// Close Ethereum client (if it has a close method)
	if ti.ETHClient != nil {
		ti.ETHClient.Close()
	}

	return nil
}
