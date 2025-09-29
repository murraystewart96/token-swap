package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	priceKeyNameSpaceFmt = "price:%s"
	priceCacheTTL        = 5 * time.Minute

	reservesKeyNameSpaceFmt = "reserves:%s"
	reservesCacheTTL        = 5 * time.Minute
)

type Cache struct {
	client *redis.Client
}

func NewCache(cfg *config.Redis) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Cache{
		client: client,
	}
}

func (c *Cache) SetPrice(ctx context.Context, pair, price string) error {
	key := fmt.Sprintf(priceKeyNameSpaceFmt, pair)

	err := c.client.Set(ctx, key, price, priceCacheTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set price: %w", err)
	}

	return nil
}

func (c *Cache) GetPrice(ctx context.Context, pair string) (string, error) {
	key := fmt.Sprintf(priceKeyNameSpaceFmt, pair)

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("price not found for pair %s: %w", pair, err)
		}
		return "", fmt.Errorf("failed to get price: %w", err)
	}

	return value, nil
}

func (c *Cache) SetReserves(ctx context.Context, poolAddr string, reserves *models.PoolReserves) error {
	key := fmt.Sprintf(reservesKeyNameSpaceFmt, poolAddr)

	reservesJSON, err := json.Marshal(reserves)
	if err != nil {
		return fmt.Errorf("failed to marshal reserves: %w", err)
	}

	err = c.client.Set(ctx, key, reservesJSON, reservesCacheTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set reserves: %w", err)
	}

	return nil
}

func (c *Cache) GetReserves(ctx context.Context, poolAddr string) (*models.PoolReserves, error) {
	key := fmt.Sprintf(reservesKeyNameSpaceFmt, poolAddr)

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("reserves not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get reserves: %w", err)
	}

	var reserves models.PoolReserves
	err = json.Unmarshal([]byte(value), &reserves)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal reserves: %w", err)
	}

	return &reserves, nil
}

func (c *Cache) Reset() error {
	return c.client.FlushDB(context.Background()).Err()
}
