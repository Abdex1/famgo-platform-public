// services/user-service/internal/infrastructure/redis_cache.go
// Redis User Cache Implementation

package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// RedisUserCache implements UserCache
type RedisUserCache struct {
	redis *redis.Client
	ttl   int32
}

func NewRedisUserCache(client *redis.Client, defaultTTL int32) *RedisUserCache {
	return &RedisUserCache{
		redis: client,
		ttl:   defaultTTL,
	}
}

// ===== USER CACHE =====

func (c *RedisUserCache) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", userID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *RedisUserCache) SetUser(ctx context.Context, user *domain.User, ttl int32) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisUserCache) DeleteUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s", userID)
	return c.redis.Del(ctx, key).Err()
}

// ===== DRIVER PROFILE CACHE =====

func (c *RedisUserCache) GetDriverProfile(ctx context.Context, driverID string) (*domain.DriverProfile, error) {
	key := fmt.Sprintf("driver:profile:%s", driverID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var profile domain.DriverProfile
	if err := json.Unmarshal([]byte(val), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *RedisUserCache) SetDriverProfile(ctx context.Context, profile *domain.DriverProfile, ttl int32) error {
	key := fmt.Sprintf("driver:profile:%s", profile.ID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisUserCache) DeleteDriverProfile(ctx context.Context, driverID string) error {
	key := fmt.Sprintf("driver:profile:%s", driverID)
	return c.redis.Del(ctx, key).Err()
}

// ===== PASSENGER PROFILE CACHE =====

func (c *RedisUserCache) GetPassengerProfile(ctx context.Context, passengerID string) (*domain.PassengerProfile, error) {
	key := fmt.Sprintf("passenger:profile:%s", passengerID)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var profile domain.PassengerProfile
	if err := json.Unmarshal([]byte(val), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *RedisUserCache) SetPassengerProfile(ctx context.Context, profile *domain.PassengerProfile, ttl int32) error {
	key := fmt.Sprintf("passenger:profile:%s", profile.ID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	duration := time.Duration(ttl) * time.Second
	if ttl == 0 {
		duration = time.Duration(c.ttl) * time.Second
	}

	return c.redis.Set(ctx, key, data, duration).Err()
}

func (c *RedisUserCache) DeletePassengerProfile(ctx context.Context, passengerID string) error {
	key := fmt.Sprintf("passenger:profile:%s", passengerID)
	return c.redis.Del(ctx, key).Err()
}
