//platform/cache/redis/config.go

package redis

import "time"

type Config struct {
	Address            string
	Password           string
	DB                 int
	PoolSize           int
	MinIdleConnections int
	MaxRetries         int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	EnableTLS          bool
}
