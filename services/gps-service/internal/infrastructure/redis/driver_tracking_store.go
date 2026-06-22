// services/gps-service/internal/infrastructure/redis/driver_tracking_store.go
// Redis store for driver online/offline status tracking

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// DriverTrackingStore handles driver status and tracking in Redis
type DriverTrackingStore struct {
	client    *redis.Client
	keyPrefix string
	ttl       time.Duration
}

// NewDriverTrackingStore creates a new driver tracking store
func NewDriverTrackingStore(client *redis.Client, keyPrefix string, ttl time.Duration) *DriverTrackingStore {
	return &DriverTrackingStore{
		client:    client,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

// StoredDriverStatus represents driver status in Redis
type StoredDriverStatus struct {
	DriverID         string    `json:"driver_id"`
	Status           string    `json:"status"`
	IsOnline         bool      `json:"is_online"`
	LastUpdateAt     time.Time `json:"last_update_at"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	AcceptanceRate   float64   `json:"acceptance_rate"`
	Rating           float64   `json:"rating"`
	RidesCompleted   int       `json:"rides_completed"`
	ConsecutiveFails int       `json:"consecutive_fails"`
}

// SetDriverOnline marks driver as online
func (d *DriverTrackingStore) SetDriverOnline(
	ctx context.Context,
	driverID string,
	location *valueobjects.Geolocation,
) error {
	if driverID == "" || location == nil {
		return fmt.Errorf("invalid driver ID or location")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	data := StoredDriverStatus{
		DriverID:     driverID,
		Status:       "online",
		IsOnline:     true,
		LastUpdateAt: time.Now(),
		Latitude:     location.Coordinates.Latitude,
		Longitude:    location.Coordinates.Longitude,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	cmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to set driver online: %w", err)
	}

	// Add to online drivers set
	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)
	d.client.SAdd(ctx, setKey, driverID)
	d.client.Expire(ctx, setKey, d.ttl)

	return nil
}

// SetDriverOffline marks driver as offline
func (d *DriverTrackingStore) SetDriverOffline(ctx context.Context, driverID string) error {
	if driverID == "" {
		return fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	// Get current status
	getCmd := d.client.Get(ctx, key)
	valStr, err := getCmd.Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err == nil {
		json.Unmarshal([]byte(valStr), &data)
	}

	data.DriverID = driverID
	data.Status = "offline"
	data.IsOnline = false
	data.LastUpdateAt = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	cmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to set driver offline: %w", err)
	}

	// Remove from online drivers set
	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)
	d.client.SRem(ctx, setKey, driverID)

	return nil
}

// GetDriverStatus retrieves driver status
func (d *DriverTrackingStore) GetDriverStatus(ctx context.Context, driverID string) (*StoredDriverStatus, error) {
	if driverID == "" {
		return nil, fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	cmd := d.client.Get(ctx, key)
	valStr, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("driver status not found")
		}
		return nil, fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err := json.Unmarshal([]byte(valStr), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver status: %w", err)
	}

	return &data, nil
}

// UpdateDriverLocation updates driver location in tracking store
func (d *DriverTrackingStore) UpdateDriverLocation(
	ctx context.Context,
	driverID string,
	location *valueobjects.Geolocation,
) error {
	if driverID == "" || location == nil {
		return fmt.Errorf("invalid driver ID or location")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	// Get current status
	cmd := d.client.Get(ctx, key)
	valStr, err := cmd.Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err == nil {
		json.Unmarshal([]byte(valStr), &data)
	}

	data.DriverID = driverID
	data.LastUpdateAt = time.Now()
	data.Latitude = location.Coordinates.Latitude
	data.Longitude = location.Coordinates.Longitude

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	setCmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := setCmd.Err(); err != nil {
		return fmt.Errorf("failed to update driver location: %w", err)
	}

	return nil
}

// GetOnlineDrivers retrieves all online drivers
func (d *DriverTrackingStore) GetOnlineDrivers(ctx context.Context) ([]string, error) {
	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)

	cmd := d.client.SMembers(ctx, setKey)
	drivers, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get online drivers: %w", err)
	}

	return drivers, nil
}

// IsDriverOnline checks if driver is online
func (d *DriverTrackingStore) IsDriverOnline(ctx context.Context, driverID string) (bool, error) {
	if driverID == "" {
		return false, fmt.Errorf("invalid driver ID")
	}

	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)

	cmd := d.client.SIsMember(ctx, setKey, driverID)
	isMember, err := cmd.Result()
	if err != nil {
		return false, fmt.Errorf("failed to check driver online status: %w", err)
	}

	return isMember, nil
}

// RecordFailure records a location update failure
func (d *DriverTrackingStore) RecordFailure(ctx context.Context, driverID string) error {
	if driverID == "" {
		return fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	cmd := d.client.Get(ctx, key)
	valStr, err := cmd.Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err == nil {
		json.Unmarshal([]byte(valStr), &data)
	}

	data.DriverID = driverID
	data.ConsecutiveFails++
	data.LastUpdateAt = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	setCmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := setCmd.Err(); err != nil {
		return fmt.Errorf("failed to record failure: %w", err)
	}

	return nil
}

// ResetFailure resets consecutive failures counter
func (d *DriverTrackingStore) ResetFailure(ctx context.Context, driverID string) error {
	if driverID == "" {
		return fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	cmd := d.client.Get(ctx, key)
	valStr, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err := json.Unmarshal([]byte(valStr), &data); err != nil {
		return fmt.Errorf("failed to unmarshal driver status: %w", err)
	}

	data.ConsecutiveFails = 0
	data.LastUpdateAt = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	setCmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := setCmd.Err(); err != nil {
		return fmt.Errorf("failed to reset failure: %w", err)
	}

	return nil
}

// UpdateStats updates driver statistics
func (d *DriverTrackingStore) UpdateStats(
	ctx context.Context,
	driverID string,
	acceptanceRate, rating float64,
	ridesCompleted int,
) error {
	if driverID == "" {
		return fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	cmd := d.client.Get(ctx, key)
	valStr, err := cmd.Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get driver status: %w", err)
	}

	var data StoredDriverStatus
	if err == nil {
		json.Unmarshal([]byte(valStr), &data)
	}

	data.DriverID = driverID
	data.AcceptanceRate = acceptanceRate
	data.Rating = rating
	data.RidesCompleted = ridesCompleted
	data.LastUpdateAt = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal driver status: %w", err)
	}

	setCmd := d.client.Set(ctx, key, jsonData, d.ttl)
	if err := setCmd.Err(); err != nil {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	return nil
}

// DeleteDriverStatus deletes driver status
func (d *DriverTrackingStore) DeleteDriverStatus(ctx context.Context, driverID string) error {
	if driverID == "" {
		return fmt.Errorf("invalid driver ID")
	}

	key := fmt.Sprintf("%s:status:%s", d.keyPrefix, driverID)

	cmd := d.client.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to delete driver status: %w", err)
	}

	// Remove from online drivers set
	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)
	d.client.SRem(ctx, setKey, driverID)

	return nil
}

// ClearAll clears all tracking data
func (d *DriverTrackingStore) ClearAll(ctx context.Context) error {
	// Get all keys with our prefix
	pattern := fmt.Sprintf("%s:status:*", d.keyPrefix)
	iter := d.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		delCmd := d.client.Del(ctx, keys...)
		if err := delCmd.Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	// Clear online drivers set
	setKey := fmt.Sprintf("%s:online_drivers", d.keyPrefix)
	delCmd := d.client.Del(ctx, setKey)
	if err := delCmd.Err(); err != nil {
		return fmt.Errorf("failed to clear online drivers: %w", err)
	}

	return nil
}
