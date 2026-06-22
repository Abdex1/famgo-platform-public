package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ActiveRidesCache caches active rides for fast O(1) lookups
type ActiveRidesCache struct {
	client      *redis.Client
	keyPrefix   string
	ttl         time.Duration
}

func NewActiveRidesCache(client *redis.Client, keyPrefix string, ttl time.Duration) *ActiveRidesCache {
	return &ActiveRidesCache{
		client:    client,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

// CachedRide represents cached ride data
type CachedRide struct {
	ID                string  `json:"id"`
	RiderID           string  `json:"rider_id"`
	DriverID          string  `json:"driver_id"`
	Status            string  `json:"status"`
	PickupLat         float64 `json:"pickup_lat"`
	PickupLng         float64 `json:"pickup_lng"`
	DropoffLat        float64 `json:"dropoff_lat"`
	DropoffLng        float64 `json:"dropoff_lng"`
	EstimatedDistance float64 `json:"estimated_distance"`
}

// Set caches a ride with TTL
func (arc *ActiveRidesCache) Set(ctx context.Context, rideID string, ride *CachedRide) error {
	if ctx == nil {
		ctx = context.Background()
	}

	key := fmt.Sprintf("%s:%s", arc.keyPrefix, rideID)
	jsonData, err := json.Marshal(ride)
	if err != nil {
		return fmt.Errorf("failed to marshal ride: %w", err)
	}

	if err := arc.client.Set(ctx, key, string(jsonData), arc.ttl).Err(); err != nil {
		return fmt.Errorf("failed to cache ride: %w", err)
	}

	// Also add to active rides set for bulk operations
	if err := arc.client.SAdd(ctx, fmt.Sprintf("%s:active", arc.keyPrefix), rideID).Err(); err != nil {
		return fmt.Errorf("failed to add to active set: %w", err)
	}

	return nil
}

// Get retrieves a cached ride
func (arc *ActiveRidesCache) Get(ctx context.Context, rideID string) (*CachedRide, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	key := fmt.Sprintf("%s:%s", arc.keyPrefix, rideID)
	jsonData := arc.client.Get(ctx, key).Val()
	if jsonData == "" {
		return nil, fmt.Errorf("ride not found in cache")
	}

	var ride CachedRide
	if err := json.Unmarshal([]byte(jsonData), &ride); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ride: %w", err)
	}

	return &ride, nil
}

// Delete removes a ride from cache
func (arc *ActiveRidesCache) Delete(ctx context.Context, rideID string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	key := fmt.Sprintf("%s:%s", arc.keyPrefix, rideID)
	if err := arc.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete ride from cache: %w", err)
	}

	if err := arc.client.SRem(ctx, fmt.Sprintf("%s:active", arc.keyPrefix), rideID).Err(); err != nil {
		return fmt.Errorf("failed to remove from active set: %w", err)
	}

	return nil
}

// GetAll retrieves all active ride IDs
func (arc *ActiveRidesCache) GetAll(ctx context.Context) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	rideIDs := arc.client.SMembers(ctx, fmt.Sprintf("%s:active", arc.keyPrefix)).Val()
	return rideIDs, nil
}

// UpdateStatus updates only the status field in cache
func (arc *ActiveRidesCache) UpdateStatus(ctx context.Context, rideID string, status string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	key := fmt.Sprintf("%s:%s", arc.keyPrefix, rideID)

	// Get existing ride
	jsonData := arc.client.Get(ctx, key).Val()
	if jsonData == "" {
		return fmt.Errorf("ride not found in cache")
	}

	var ride CachedRide
	if err := json.Unmarshal([]byte(jsonData), &ride); err != nil {
		return fmt.Errorf("failed to unmarshal ride: %w", err)
	}

	// Update status
	ride.Status = status
	jsonData, err := json.Marshal(ride)
	if err != nil {
		return fmt.Errorf("failed to marshal ride: %w", err)
	}

	if err := arc.client.Set(ctx, key, string(jsonData), arc.ttl).Err(); err != nil {
		return fmt.Errorf("failed to update cache: %w", err)
	}

	return nil
}

// Clear removes all rides from cache
func (arc *ActiveRidesCache) Clear(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	keys := arc.client.Keys(ctx, fmt.Sprintf("%s:*", arc.keyPrefix)).Val()
	if len(keys) > 0 {
		if err := arc.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	return nil
}
