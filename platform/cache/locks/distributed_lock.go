
//platform/cache/locks/distributed_lock.go

package locks

import (
    "context"
    "time"

    goredis "github.com/redis/go-redis/v9"
)

type DriverLockService struct {
    rdb *goredis.Client
}

func NewDriverLockService(
    rdb *goredis.Client,
) *DriverLockService {
    return &DriverLockService{
        rdb: rdb,
    }
}

func (d *DriverLockService) Acquire(
    ctx context.Context,
    driverID string,
    ttl time.Duration,
) (bool, error) {

    return d.rdb.SetNX(
        ctx,
        "lock:driver:"+driverID,
        "1",
        ttl,
    ).Result()
}

func (d *DriverLockService) Release(
    ctx context.Context,
    driverID string,
) error {

    return d.rdb.Del(
        ctx,
        "lock:driver:"+driverID,
    ).Err()
}
