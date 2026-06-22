/*
# PHASE 8 — DISTRIBUTED LOCKS

CRITICAL FOR:

* dispatch
* payments
* session rotation
* fraud

---

# STEP 11 — CREATE LOCK MANAGER

packages/redis-platform/locks/locks.go
*/
package locks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	rdb *redis.Client
}

func NewManager(rdb *redis.Client) *Manager {
	return &Manager{rdb: rdb}
}

func (m *Manager) Acquire(
	ctx context.Context,
	key string,
	ttl time.Duration,
) (bool, error) {

	return m.rdb.SetNX(
		ctx,
		"lock:"+key,
		"1",
		ttl,
	).Result()
}

func (m *Manager) Release(
	ctx context.Context,
	key string,
) error {

	return m.rdb.Del(
		ctx,
		"lock:"+key,
	).Err()
}
