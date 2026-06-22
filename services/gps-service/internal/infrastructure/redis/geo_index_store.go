// services/gps-service/internal/infrastructure/redis/geo_index_store.go
// Redis GEO index store for fast spatial queries

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// GeoIndexStore handles Redis GEO operations
type GeoIndexStore struct {
	client    *redis.Client
	keyPrefix string
	ttl       time.Duration
}

// NewGeoIndexStore creates a new geo index store
func NewGeoIndexStore(client *redis.Client, keyPrefix string, ttl time.Duration) *GeoIndexStore {
	return &GeoIndexStore{
		client:    client,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

// AddLocation adds a driver location to the GEO index
func (g *GeoIndexStore) AddLocation(
	ctx context.Context,
	driverID string,
	location *valueobjects.Geolocation,
) error {
	if driverID == "" || location == nil || !location.IsValid() {
		return fmt.Errorf("invalid driver ID or location")
	}

	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	// GEOADD key longitude latitude member
	cmd := g.client.GeoAdd(
		ctx,
		key,
		&redis.GeoLocation{
			Name:      driverID,
			Longitude: location.Coordinates.Longitude,
			Latitude:  location.Coordinates.Latitude,
		},
	)

	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to add location to GEO index: %w", err)
	}

	// Set TTL for the entire key
	g.client.Expire(ctx, key, g.ttl)

	return nil
}

// AddLocations adds multiple driver locations to the GEO index
func (g *GeoIndexStore) AddLocations(
	ctx context.Context,
	locations map[string]*valueobjects.Geolocation,
) error {
	if len(locations) == 0 {
		return nil
	}

	key := fmt.Sprintf("%s:drivers", g.keyPrefix)
	geoLocations := make([]*redis.GeoLocation, 0, len(locations))

	for driverID, location := range locations {
		if location != nil && location.IsValid() {
			geoLocations = append(geoLocations, &redis.GeoLocation{
				Name:      driverID,
				Longitude: location.Coordinates.Longitude,
				Latitude:  location.Coordinates.Latitude,
			})
		}
	}

	if len(geoLocations) == 0 {
		return nil
	}

	cmd := g.client.GeoAdd(ctx, key, geoLocations...)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to add locations to GEO index: %w", err)
	}

	g.client.Expire(ctx, key, g.ttl)
	return nil
}

// FindNearby finds drivers within radius of a location
func (g *GeoIndexStore) FindNearby(
	ctx context.Context,
	centerLocation *valueobjects.Geolocation,
	radiusKm float64,
	limit int,
) ([]string, error) {
	if centerLocation == nil || !centerLocation.IsValid() {
		return nil, fmt.Errorf("invalid center location")
	}

	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	// GEORADIUS key longitude latitude radius km COUNT limit WITHDIST ASC
	cmd := g.client.GeoRadius(
		ctx,
		key,
		centerLocation.Coordinates.Longitude,
		centerLocation.Coordinates.Latitude,
		&redis.GeoRadiusQuery{
			Radius:      radiusKm,
			Unit:        "km",
			Count:       int64(limit),
			Sort:        "ASC",
			WithCoord:   true,
			WithDist:    true,
		},
	)

	results, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to find nearby drivers: %w", err)
	}

	driverIDs := make([]string, 0, len(results))
	for _, result := range results {
		driverIDs = append(driverIDs, result.Name)
	}

	return driverIDs, nil
}

// FindNearbyWithDistance finds drivers and returns distances
func (g *GeoIndexStore) FindNearbyWithDistance(
	ctx context.Context,
	centerLocation *valueobjects.Geolocation,
	radiusKm float64,
	limit int,
) ([]map[string]interface{}, error) {
	if centerLocation == nil || !centerLocation.IsValid() {
		return nil, fmt.Errorf("invalid center location")
	}

	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.GeoRadius(
		ctx,
		key,
		centerLocation.Coordinates.Longitude,
		centerLocation.Coordinates.Latitude,
		&redis.GeoRadiusQuery{
			Radius:      radiusKm,
			Unit:        "km",
			Count:       int64(limit),
			Sort:        "ASC",
			WithCoord:   true,
			WithDist:    true,
		},
	)

	results, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("failed to find nearby drivers with distance: %w", err)
	}

	var driverData []map[string]interface{}
	for _, result := range results {
		driverData = append(driverData, map[string]interface{}{
			"driver_id": result.Name,
			"distance":  result.Dist,
			"latitude":  result.Longitude,
			"longitude": result.Latitude,
		})
	}

	return driverData, nil
}

// GetDistance calculates distance between two members
func (g *GeoIndexStore) GetDistance(
	ctx context.Context,
	member1, member2 string,
) (float64, error) {
	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.GeoDist(ctx, key, member1, member2, "km")
	distance, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("one or both members not found")
		}
		return 0, fmt.Errorf("failed to get distance: %w", err)
	}

	return distance, nil
}

// GetCoordinates gets coordinates for a member
func (g *GeoIndexStore) GetCoordinates(
	ctx context.Context,
	driverID string,
) (*valueobjects.Coordinates, error) {
	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.GeoPos(ctx, key, driverID)
	results, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates: %w", err)
	}

	if len(results) == 0 || results[0] == nil {
		return nil, fmt.Errorf("driver not found in GEO index")
	}

	coords, err := valueobjects.NewCoordinates(results[0].Latitude, results[0].Longitude)
	if err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	return coords, nil
}

// RemoveLocation removes a driver from the GEO index
func (g *GeoIndexStore) RemoveLocation(ctx context.Context, driverID string) error {
	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.ZRem(ctx, key, driverID)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to remove location: %w", err)
	}

	return nil
}

// RemoveLocations removes multiple drivers from the GEO index
func (g *GeoIndexStore) RemoveLocations(ctx context.Context, driverIDs ...string) error {
	if len(driverIDs) == 0 {
		return nil
	}

	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.ZRem(ctx, key, convertToInterface(driverIDs)...)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to remove locations: %w", err)
	}

	return nil
}

// ClearAll clears all locations from the GEO index
func (g *GeoIndexStore) ClearAll(ctx context.Context) error {
	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to clear GEO index: %w", err)
	}

	return nil
}

// GetCount returns the number of drivers in the GEO index
func (g *GeoIndexStore) GetCount(ctx context.Context) (int64, error) {
	key := fmt.Sprintf("%s:drivers", g.keyPrefix)

	cmd := g.client.ZCard(ctx, key)
	count, err := cmd.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get count: %w", err)
	}

	return count, nil
}

// UpdateLocationInBatch updates multiple locations efficiently
func (g *GeoIndexStore) UpdateLocationInBatch(
	ctx context.Context,
	drivers []*entities.DriverLocation,
) error {
	locations := make(map[string]*valueobjects.Geolocation)
	for _, driver := range drivers {
		if driver != nil && driver.CurrentLocation != nil {
			locations[driver.DriverID] = driver.CurrentLocation
		}
	}

	return g.AddLocations(ctx, locations)
}

// GetNearbyDriversForPickup finds drivers near a pickup location
func (g *GeoIndexStore) GetNearbyDriversForPickup(
	ctx context.Context,
	pickupLocation *valueobjects.Geolocation,
	radiusKm float64,
	limit int,
) ([]string, error) {
	return g.FindNearby(ctx, pickupLocation, radiusKm, limit)
}

// ExpandSearchIfNeeded expands search radius if not enough drivers found
func (g *GeoIndexStore) ExpandSearchIfNeeded(
	ctx context.Context,
	centerLocation *valueobjects.Geolocation,
	initialRadius float64,
	minDrivers int,
	maxRadius float64,
) ([]string, float64, error) {
	radius := initialRadius
	geoService := services.NewRedisGeoService()

	for radius <= maxRadius {
		drivers, err := g.FindNearby(ctx, centerLocation, radius, 1000)
		if err != nil {
			return nil, 0, err
		}

		if len(drivers) >= minDrivers {
			return drivers, radius, nil
		}

		// Expand radius
		radius = geoService.ExpandSearchArea(radius, len(drivers), minDrivers, maxRadius)
		if radius == initialRadius {
			break
		}
	}

	// Return all found drivers even if less than minimum
	drivers, err := g.FindNearby(ctx, centerLocation, maxRadius, 1000)
	return drivers, maxRadius, err
}

// Helper function to convert string slice to interface slice
func convertToInterface(s []string) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}
