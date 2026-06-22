// services/ride-service/internal/infrastructure/redis_cache.go
// Redis Ride Cache Implementation

package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/packages/redis-platform"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// RedisRideCache implements RideCache using packages/redis-platform (Rule 2 compliant)
type RedisRideCache struct {
	redis redis_platform.RedisClient  // Using platform wrapper, not raw redis
	ttl   int32
}

func NewRedisRideCache(client redis_platform.RedisClient, defaultTTL int32) *RedisRideCache {
	return &RedisRideCache{
		redis: client,
		ttl:   defaultTTL,
	}
}

// ===== RIDE CACHE METHODS =====

func (c *RedisRideCache) GetRide(ctx context.Context, rideID string) (*domain.Ride, error) {
	key := fmt.Sprintf("ride:%s", rideID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var ride domain.Ride
	if err := json.Unmarshal([]byte(val), &ride); err != nil {
		return nil, err
	}

	return &ride, nil
}

func (c *RedisRideCache) SetRide(ctx context.Context, ride *domain.Ride, ttl int32) error {
	key := fmt.Sprintf("ride:%s", ride.ID)
	data, err := json.Marshal(ride)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisRideCache) DeleteRide(ctx context.Context, rideID string) error {
	key := fmt.Sprintf("ride:%s", rideID)
	return c.redis.Del(ctx, key).Err()
}

// ===== ACTIVE RIDES CACHE =====

func (c *RedisRideCache) GetActiveRides(ctx context.Context) ([]domain.Ride, error) {
	key := "rides:active"
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return []domain.Ride{}, nil
		}
		return nil, err
	}

	var rides []domain.Ride
	if err := json.Unmarshal([]byte(val), &rides); err != nil {
		return nil, err
	}

	return rides, nil
}

func (c *RedisRideCache) SetActiveRides(ctx context.Context, rides []domain.Ride, ttl int32) error {
	key := "rides:active"
	data, err := json.Marshal(rides)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

// ===== PASSENGER ACTIVE RIDES CACHE =====

func (c *RedisRideCache) GetPassengerActiveRides(ctx context.Context, passengerID string) ([]domain.Ride, error) {
	key := fmt.Sprintf("rides:passenger:active:%s", passengerID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return []domain.Ride{}, nil
		}
		return nil, err
	}

	var rides []domain.Ride
	if err := json.Unmarshal([]byte(val), &rides); err != nil {
		return nil, err
	}

	return rides, nil
}

func (c *RedisRideCache) SetPassengerActiveRides(ctx context.Context, passengerID string, rides []domain.Ride, ttl int32) error {
	key := fmt.Sprintf("rides:passenger:active:%s", passengerID)
	data, err := json.Marshal(rides)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisRideCache) DeletePassengerActiveRides(ctx context.Context, passengerID string) error {
	key := fmt.Sprintf("rides:passenger:active:%s", passengerID)
	return c.redis.Del(ctx, key).Err()
}

// ===== DRIVER ACTIVE RIDES CACHE =====

func (c *RedisRideCache) GetDriverActiveRides(ctx context.Context, driverID string) ([]domain.Ride, error) {
	key := fmt.Sprintf("rides:driver:active:%s", driverID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return []domain.Ride{}, nil
		}
		return nil, err
	}

	var rides []domain.Ride
	if err := json.Unmarshal([]byte(val), &rides); err != nil {
		return nil, err
	}

	return rides, nil
}

func (c *RedisRideCache) SetDriverActiveRides(ctx context.Context, driverID string, rides []domain.Ride, ttl int32) error {
	key := fmt.Sprintf("rides:driver:active:%s", driverID)
	data, err := json.Marshal(rides)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisRideCache) DeleteDriverActiveRides(ctx context.Context, driverID string) error {
	key := fmt.Sprintf("rides:driver:active:%s", driverID)
	return c.redis.Del(ctx, key).Err()
}

// ===== RIDE SEARCH CACHE (for dispatch service integration) =====

func (c *RedisRideCache) GetRidesNearby(ctx context.Context, lat, lon float64, radiusKm int) ([]domain.Ride, error) {
	// GEORADIUS implementation would go here
	// For now, return from general active rides
	return c.GetActiveRides(ctx)
}

// DeleteAllCache invalidates all ride-related caches
func (c *RedisRideCache) DeleteAllCache(ctx context.Context) error {
	pattern := "rides:*"
	iter := c.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
