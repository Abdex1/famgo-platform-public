package redis

import (
	"context"
	"fmt"
	"strings"

	redisclient "github.com/Abdex1/FamGo-platform/packages/redis-platform/client"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewFromURL(redisURL string, db int) (*Client, error) {
	address := strings.TrimPrefix(redisURL, "redis://")
	if address == "" {
		address = "localhost:6379"
	}
	rdb := redisclient.New(redisclient.Config{
		Address: address,
		DB:      db,
	})
	ctx := context.Background()
	if err := redisclient.HealthCheck(ctx, rdb); err != nil {
		return nil, fmt.Errorf("redis health check failed: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

func (c *Client) Raw() *redis.Client {
	return c.rdb
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
