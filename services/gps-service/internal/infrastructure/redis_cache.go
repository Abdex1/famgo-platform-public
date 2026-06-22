package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/packages/redis-platform"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
)

// RedisLocationCache provides caching for driver locations
type RedisLocationCache struct {
	redis redis_platform.RedisClient
	ttl   time.Duration
}

// NewRedisLocationCache creates a new Redis location cache
func NewRedisLocationCache(redis redis_platform.RedisClient, ttl time.Duration) *RedisLocationCache {
	return &RedisLocationCache{
		redis: redis,
		ttl:   ttl,
	}
}

// GetDriverLocation retrieves cached location for a driver
// Key format: gps:location:{driver_id}
func (c *RedisLocationCache) GetDriverLocation(ctx context.Context, driverID string) (*domain.DriverLocation, error) {
	key := fmt.Sprintf("gps:location:%s", driverID)

	val, err := c.redis.Get(ctx, key)
	if err != nil {
		// Cache miss - return nil (not an error)
		return nil, nil
	}

	var loc domain.DriverLocation
	err = json.Unmarshal(val, &loc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached location: %w", err)
	}

	return &loc, nil
}

// CacheDriverLocation stores a driver location in Redis
func (c *RedisLocationCache) CacheDriverLocation(ctx context.Context, location *domain.DriverLocation) error {
	key := fmt.Sprintf("gps:location:%s", location.DriverID)

	data, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	err = c.redis.SetEX(ctx, key, data, c.ttl)
	if err != nil {
		return fmt.Errorf("failed to cache location: %w", err)
	}

	return nil
}

// InvalidateDriverLocation removes cached location
func (c *RedisLocationCache) InvalidateDriverLocation(ctx context.Context, driverID string) error {
	key := fmt.Sprintf("gps:location:%s", driverID)

	err := c.redis.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to invalidate location: %w", err)
	}

	return nil
}

// ========================================

// RedisTripCache provides caching for trips
type RedisTripCache struct {
	redis redis_platform.RedisClient
	ttl   time.Duration
}

// NewRedisTripCache creates a new Redis trip cache
func NewRedisTripCache(redis redis_platform.RedisClient, ttl time.Duration) *RedisTripCache {
	return &RedisTripCache{
		redis: redis,
		ttl:   ttl,
	}
}

// GetTrip retrieves cached trip
// Key format: gps:trip:{trip_id}
func (c *RedisTripCache) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	key := fmt.Sprintf("gps:trip:%s", tripID)

	val, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil
	}

	var trip domain.Trip
	err = json.Unmarshal(val, &trip)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached trip: %w", err)
	}

	return &trip, nil
}

// CacheTrip stores a trip in Redis
func (c *RedisTripCache) CacheTrip(ctx context.Context, trip *domain.Trip) error {
	key := fmt.Sprintf("gps:trip:%s", trip.ID)

	data, err := json.Marshal(trip)
	if err != nil {
		return fmt.Errorf("failed to marshal trip: %w", err)
	}

	err = c.redis.SetEX(ctx, key, data, c.ttl)
	if err != nil {
		return fmt.Errorf("failed to cache trip: %w", err)
	}

	return nil
}

// ========================================

// RedisGeofenceCache provides caching for geofences
type RedisGeofenceCache struct {
	redis redis_platform.RedisClient
	ttl   time.Duration
}

// NewRedisGeofenceCache creates a new Redis geofence cache
func NewRedisGeofenceCache(redis redis_platform.RedisClient, ttl time.Duration) *RedisGeofenceCache {
	return &RedisGeofenceCache{
		redis: redis,
		ttl:   ttl,
	}
}

// GetAllGeofences retrieves cached geofences
// Key format: gps:geofences:all
func (c *RedisGeofenceCache) GetAllGeofences(ctx context.Context) ([]domain.Geofence, error) {
	key := "gps:geofences:all"

	val, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil
	}

	var geofences []domain.Geofence
	err = json.Unmarshal(val, &geofences)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached geofences: %w", err)
	}

	return geofences, nil
}

// CacheAllGeofences stores all geofences in Redis
func (c *RedisGeofenceCache) CacheAllGeofences(ctx context.Context, geofences []domain.Geofence) error {
	key := "gps:geofences:all"

	data, err := json.Marshal(geofences)
	if err != nil {
		return fmt.Errorf("failed to marshal geofences: %w", err)
	}

	err = c.redis.SetEX(ctx, key, data, c.ttl)
	if err != nil {
		return fmt.Errorf("failed to cache geofences: %w", err)
	}

	return nil
}

// InvalidateGeofences removes cached geofences
func (c *RedisGeofenceCache) InvalidateGeofences(ctx context.Context) error {
	key := "gps:geofences:all"

	err := c.redis.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to invalidate geofences: %w", err)
	}

	return nil
}

// ========================================

// RedisDriverCache provides caching for active drivers using geospatial queries
type RedisDriverCache struct {
	redis redis_platform.RedisClient
}

// NewRedisDriverCache creates a new Redis driver cache
func NewRedisDriverCache(redis redis_platform.RedisClient) *RedisDriverCache {
	return &RedisDriverCache{redis: redis}
}

// AddDriver adds a driver to the active drivers cache
// Key format: gps:drivers:geo (Redis GEOADD)
func (c *RedisDriverCache) AddDriver(ctx context.Context, driverID string, latitude float64, longitude float64) error {
	key := "gps:drivers:geo"

	err := c.redis.GeoAdd(ctx, key, latitude, longitude, driverID)
	if err != nil {
		return fmt.Errorf("failed to add driver to geo cache: %w", err)
	}

	return nil
}

// GetNearbyDrivers retrieves drivers within a radius
// Radius in meters
func (c *RedisDriverCache) GetNearbyDrivers(ctx context.Context, latitude float64, longitude float64, radiusMeters float64) ([]string, error) {
	key := "gps:drivers:geo"

	// Convert radius from meters to km for Redis GEORADIUS
	radiusKm := radiusMeters / 1000.0

	drivers, err := c.redis.GeoRadius(ctx, key, latitude, longitude, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby drivers: %w", err)
	}

	return drivers, nil
}

// RemoveDriver removes a driver from the active drivers cache
func (c *RedisDriverCache) RemoveDriver(ctx context.Context, driverID string) error {
	key := "gps:drivers:geo"

	err := c.redis.Delete(ctx, driverID)
	if err != nil {
		return fmt.Errorf("failed to remove driver from geo cache: %w", err)
	}

	return nil
}
