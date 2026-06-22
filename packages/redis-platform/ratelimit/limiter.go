/*
# PHASE 10 — RATE LIMITING

---

# STEP 13 — CREATE RATE LIMITER

packages/redis-platform/ratelimit/limiter.go
*/
package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb *redis.Client
}

func NewLimiter(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb}
}

func (l *Limiter) Allow(
	ctx context.Context,
	key string,
	limit int,
	window time.Duration,
) (bool, error) {

	redisKey := fmt.Sprintf("ratelimit:%s", key)

	count, err := l.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		_ = l.rdb.Expire(ctx, redisKey, window).Err()
	}

	return count <= int64(limit), nil
}
