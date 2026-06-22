package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain"
	"github.com/go-redis/redis/v8"
)

// LocationStore handles location storage in Redis
type LocationStore struct {
	client *redis.Client
}

// NewLocationStore creates new location store
func NewLocationStore(client *redis.Client) *LocationStore {
	return &LocationStore{client: client}
}

// SaveDriverLocation saves current driver location to Redis GEO index and hash
func (s *LocationStore) SaveDriverLocation(ctx context.Context, loc *domain.DriverLocation) error {
	// Store in GEO index for spatial queries
	cmd := s.client.GeoAdd(ctx, "drivers:geo", &redis.GeoLocation{
		Name:      loc.DriverID,
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
	})

	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to add driver to geo index: %w", err)
	}

	// Store location details in hash
	locData, _ := json.Marshal(loc)
	return s.client.HSet(ctx, "drivers:locations", loc.DriverID, string(locData)).Err()
}

// GetDriverLocation retrieves driver's current location
func (s *LocationStore) GetDriverLocation(ctx context.Context, driverID string) (*domain.DriverLocation, error) {
	data, err := s.client.HGet(ctx, "drivers:locations", driverID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("driver location not found")
	}
	if err != nil {
		return nil, err
	}

	var loc domain.DriverLocation
	if err := json.Unmarshal([]byte(data), &loc); err != nil {
		return nil, err
	}

	return &loc, nil
}

// GetNearbyDrivers finds drivers within radius using GEO index
func (s *LocationStore) GetNearbyDrivers(ctx context.Context, lat, lng, radiusMeters float64) ([]domain.NearbyDriver, error) {
	// Convert meters to km
	radiusKm := radiusMeters / 1000.0

	cmd := s.client.GeoRadius(ctx, "drivers:geo", lng, lat, &redis.GeoRadiusQuery{
		Radius: radiusKm,
		Unit:   "km",
		Count:  50,
	})

	locations, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby drivers: %w", err)
	}

	var drivers []domain.NearbyDriver
	for _, loc := range locations {
		distance := s.calculateDistance(lat, lng, loc.Latitude, loc.Longitude)
		drivers = append(drivers, domain.NearbyDriver{
			DriverID:  loc.Name,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Distance:  distance,
		})
	}

	return drivers, nil
}

// SaveLocationHistory saves location point to time series (Redis Stream)
func (s *LocationStore) SaveLocationHistory(ctx context.Context, update *domain.LocationUpdate) error {
	data, _ := json.Marshal(update)
	
	// Use Redis Stream for time-series data
	cmd := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: fmt.Sprintf("locations:%s", update.RideID),
		Values: map[string]interface{}{
			"data":      string(data),
			"timestamp": update.Timestamp.Unix(),
		},
		MaxLen: 10000, // Keep last 10k points
	})

	return cmd.Err()
}

// GetLocationHistory retrieves location history for ride
func (s *LocationStore) GetLocationHistory(ctx context.Context, rideID string) ([]domain.LocationUpdate, error) {
	cmd := s.client.XRange(ctx, fmt.Sprintf("locations:%s", rideID), "-", "+")
	
	messages, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	var updates []domain.LocationUpdate
	for _, msg := range messages {
		if data, ok := msg.Values["data"].(string); ok {
			var update domain.LocationUpdate
			if err := json.Unmarshal([]byte(data), &update); err == nil {
				updates = append(updates, update)
			}
		}
	}

	return updates, nil
}

// SetDriverOnline marks driver as online in Redis set
func (s *LocationStore) SetDriverOnline(ctx context.Context, driverID string) error {
	return s.client.SAdd(ctx, "drivers:online", driverID).Err()
}

// SetDriverOffline removes driver from online set
func (s *LocationStore) SetDriverOffline(ctx context.Context, driverID string) error {
	return s.client.SRem(ctx, "drivers:online", driverID).Err()
}

// GetOnlineDrivers retrieves all online driver IDs
func (s *LocationStore) GetOnlineDrivers(ctx context.Context) ([]string, error) {
	return s.client.SMembers(ctx, "drivers:online").Result()
}

// Helper function to calculate distance
func (s *LocationStore) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // meters
	const toRad = 3.14159265359 / 180

	dLat := (lat2 - lat1) * toRad
	dLng := (lng2 - lng1) * toRad

	a := 0.5 - cos((lat2-lat1)*toRad)/2 + cos(lat1*toRad)*cos(lat2*toRad)*(1-cos((lng2-lng1)*toRad))/2
	return 2 * R * asin(sqrt(a))
}

func cos(x float64) float64 {
	// Simplified cosine approximation
	return 1 - x*x/2 + x*x*x*x/24
}

func sin(x float64) float64 {
	// Simplified sine approximation
	return x - x*x*x/6 + x*x*x*x*x/120
}

func asin(x float64) float64 {
	// Simplified arcsin approximation
	return x + x*x*x/6 + 3*x*x*x*x*x/40
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	// Newton's method for square root
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
