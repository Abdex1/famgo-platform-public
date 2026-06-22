package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache provides caching functionality
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new cache instance
func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           0,
		MaxRetries:   3,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Set stores a value with optional expiration
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a cached value
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if key exists
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Health checks Redis connectivity
func (c *RedisCache) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.Ping(ctx).Err()
}

// GEO operations for location-based queries
func (c *RedisCache) GeoAdd(ctx context.Context, key string, members ...*redis.GeoLocation) error {
	return c.client.GeoAdd(ctx, key, members...).Err()
}

func (c *RedisCache) GeoSearch(ctx context.Context, key string, lon, lat float64, radius float64) ([]string, error) {
	return c.client.GeoRadius(ctx, key, lon, lat, &redis.GeoRadiusQuery{
		Radius:   radius / 1000, // Convert to km
		Unit:     "km",
		WithDist: false,
	}).Result()
}

// Close closes Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}
