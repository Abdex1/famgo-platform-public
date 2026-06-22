//packages/redis-platform/client/client.go

package client

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address      string
	Password     string
	DB            int
	PoolSize      int
	MinIdleConns  int
	MaxRetries    int
	TLSEnabled    bool
}

func New(cfg Config) *redis.Client {
	opts := &redis.Options{
		Addr:            cfg.Address,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxRetries:      cfg.MaxRetries,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		ConnMaxLifetime: 30 * time.Minute,
	}

	if cfg.TLSEnabled {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return redis.NewClient(opts)
}

func HealthCheck(ctx context.Context, rdb *redis.Client) error {
	return rdb.Ping(ctx).Err()
}
