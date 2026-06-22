//3. ENTERPRISE REDIS PLATFORM

//platform/cache/redis/client.go

package redis

import (
    "context"
    "time"

    goredis "github.com/redis/go-redis/v9"
)

type Config struct {
    Address            string
    Username           string
    Password           string
    Database           int
    EnableTLS          bool
    PoolSize           int
    MinIdleConnections int
}

type Client struct {
    rdb *goredis.Client
}

func New(cfg Config) *Client {
    rdb := goredis.NewClient(&goredis.Options{
        Addr:            cfg.Address,
        Username:        cfg.Username,
        Password:        cfg.Password,
        DB:              cfg.Database,
        PoolSize:        cfg.PoolSize,
        MinIdleConns:    cfg.MinIdleConnections,
        ConnMaxLifetime: 30 * time.Minute,
        DialTimeout:     10 * time.Second,
        ReadTimeout:     5 * time.Second,
        WriteTimeout:    5 * time.Second,
    })

    return &Client{rdb: rdb}
}

func (c *Client) Health(ctx context.Context) error {
    return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
    return c.rdb.Close()
}
