package security

import (
	"context"
	"time"

	"github.com/Abdex1/FamGo-platform/packages/redis-platform/ratelimit"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter adapts redis-platform ratelimit to the security middleware interface.
type RedisRateLimiter struct {
	limiter *ratelimit.Limiter
}

func NewRedisRateLimiter(rdb *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{limiter: ratelimit.NewLimiter(rdb)}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	return r.limiter.Allow(ctx, key, limit, window)
}
