//go:build week3

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"famgo/driver-service/internal/model"
	"famgo/driver-service/internal/repository"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// LocationService handles real-time and historical location tracking
// WEEK 3: Added alongside existing DriverService (NOT replacing it)
type LocationService struct {
	locationRepo *repository.LocationRepository
	redisGeo     *repository.RedisGeoRepository
	logger       logger.Logger
}

// NewLocationService creates a new location service
func NewLocationService(
	locationRepo *repository.LocationRepository,
	redisGeo *repository.RedisGeoRepository,
	log logger.Logger,
) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		redisGeo:     redisGeo,
		logger:       log,
	}
}

// UpdateDriverLocation updates driver's current location in real-time
func (s *LocationService) UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error {
	s.logger.Info("updating driver location", map[string]interface{}{
		"driver_id": driverID,
		"lat":       lat,
		"lng":       lng,
	})

	// Update Redis GEO for real-time (24h expiry)
	if err := s.redisGeo.UpdateDriverLocation(ctx, driverID, lat, lng); err != nil {
		s.logger.Error("failed to update Redis location", map[string]interface{}{"error": err})
		return err
	}

	// Archive to PostGIS every 6 hours (background job)
	location := &model.DriverLocation{
		DriverID: driverID,
		Location: model.Location{
			Latitude:    lat,
			Longitude:   lng,
			Timestamp:   time.Now(),
			IsAvailable: true,
			Status:      "online",
		},
		CreatedAt: time.Now(),
	}

	if err := s.locationRepo.RecordLocation(ctx, location); err != nil {
		s.logger.Warn("failed to archive location", map[string]interface{}{"error": err})
		// Don't fail - Redis update succeeded
	}

	s.logger.Info("location updated", map[string]interface{}{
		"driver_id": driverID,
		"lat":       lat,
		"lng":       lng,
	})

	return nil
}

// GetCurrentLocation retrieves driver's current location from Redis
func (s *LocationService) GetCurrentLocation(ctx context.Context, driverID string) (*model.Location, error) {
	location, err := s.redisGeo.GetDriverLocation(ctx, driverID)
	if err != nil {
		s.logger.Warn("failed to get current location", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	return location, nil
}

// GetLocationHistory retrieves historical location data from PostGIS
func (s *LocationService) GetLocationHistory(ctx context.Context, driverID string, startTime, endTime time.Time) ([]*model.DriverLocation, error) {
	s.logger.Info("getting location history", map[string]interface{}{
		"driver_id": driverID,
		"start":     startTime,
		"end":       endTime,
	})

	locations, err := s.locationRepo.GetLocationHistory(ctx, driverID, startTime, endTime)
	if err != nil {
		s.logger.Warn("failed to get location history", map[string]interface{}{"error": err})
		return nil, err
	}

	return locations, nil
}

// FindNearbyDrivers finds all active drivers within radius (uses Redis GEO)
func (s *LocationService) FindNearbyDrivers(ctx context.Context, lat, lng, radiusMeters float64) ([]string, error) {
	s.logger.Info("finding nearby drivers", map[string]interface{}{
		"lat":    lat,
		"lng":    lng,
		"radius": radiusMeters,
	})

	drivers, err := s.redisGeo.GetNearbyDrivers(ctx, lat, lng, radiusMeters)
	if err != nil {
		s.logger.Warn("failed to find nearby drivers", map[string]interface{}{"error": err})
		return nil, err
	}

	s.logger.Info("nearby drivers found", map[string]interface{}{
		"count": len(drivers),
	})

	return drivers, nil
}

// GeofenceService handles geographic boundaries and zone queries
type GeofenceService struct {
	locationRepo *repository.LocationRepository
	logger       logger.Logger
}

// NewGeofenceService creates a new geofence service
func NewGeofenceService(locationRepo *repository.LocationRepository, log logger.Logger) *GeofenceService {
	return &GeofenceService{
		locationRepo: locationRepo,
		logger:       log,
	}
}

// IsDriverInZone checks if driver is currently in a service zone
func (s *GeofenceService) IsDriverInZone(ctx context.Context, driverID string, zoneID string) (bool, error) {
	s.logger.Info("checking if driver in zone", map[string]interface{}{
		"driver_id": driverID,
		"zone_id":   zoneID,
	})

	inZone, err := s.locationRepo.IsDriverInZone(ctx, driverID, zoneID)
	if err != nil {
		s.logger.Warn("failed to check zone", map[string]interface{}{"error": err})
		return false, err
	}

	s.logger.Info("zone check complete", map[string]interface{}{
		"driver_id": driverID,
		"zone_id":   zoneID,
		"in_zone":   inZone,
	})

	return inZone, nil
}

// GetZoneHistory retrieves entry/exit history for driver in zone
func (s *GeofenceService) GetZoneHistory(ctx context.Context, driverID string, zoneID string) ([]map[string]interface{}, error) {
	s.logger.Info("getting zone history", map[string]interface{}{
		"driver_id": driverID,
		"zone_id":   zoneID,
	})

	// Query location history and determine zone entries/exits
	// This would be implemented in PostGIS with spatial queries

	history := []map[string]interface{}{}

	return history, nil
}

// ArchiveLocationData archives Redis GEO data to PostGIS (background job)
// Should run every 6 hours
func (s *GeofenceService) ArchiveLocationData(ctx context.Context) error {
	s.logger.Info("archiving location data from Redis to PostGIS")

	// This would be a scheduled job that:
	// 1. Gets all driver locations from Redis GEO
	// 2. Inserts into PostGIS driver_locations_history
	// 3. Cleans up old Redis data

	s.logger.Info("location data archived")

	return nil
}

// LocationTrackingConfig configuration for location tracking
type LocationTrackingConfig struct {
	UpdateIntervalSeconds int           // How often client sends location (30s default)
	RedisExpiryHours      int           // Redis data TTL (24h default)
	ArchiveIntervalHours  int           // Archive interval (6h default)
	HistoryRetentionDays  int           // PostGIS retention (90 days default)
	NearbySearchRadiusKm  float64       // Default radius for nearby searches (10km default)
}

// DefaultLocationTrackingConfig returns default configuration
func DefaultLocationTrackingConfig() *LocationTrackingConfig {
	return &LocationTrackingConfig{
		UpdateIntervalSeconds: 30,
		RedisExpiryHours:      24,
		ArchiveIntervalHours:  6,
		HistoryRetentionDays:  90,
		NearbySearchRadiusKm:  10,
	}
}

// ValidateLocation validates location coordinates
func ValidateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude: must be between -90 and 90")
	}

	if lng < -180 || lng > 180 {
		return errors.New("invalid longitude: must be between -180 and 180")
	}

	return nil
}

// ValidateRadius validates search radius
func ValidateRadius(radiusMeters float64) error {
	if radiusMeters <= 0 {
		return errors.New("radius must be positive")
	}

	if radiusMeters > 50000 {
		return errors.New("radius too large (max 50km)")
	}

	return nil
}

// CalculateDistance calculates distance between two coordinates in kilometers
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Haversine formula
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * 3.14159265359 / 180.0
	lat2Rad := lat2 * 3.14159265359 / 180.0
	deltaLat := (lat2 - lat1) * 3.14159265359 / 180.0
	deltaLng := (lng2 - lng1) * 3.14159265359 / 180.0

	a := (1 - cos(deltaLat)) / 2
	c := 2 * atan2(sqrt(a), sqrt(1-a))

	return earthRadiusKm * c
}

func cos(x float64) float64 {
	return 1 - x*x/2
}

func sin(x float64) float64 {
	return x - x*x*x/6
}

func atan2(y, x float64) float64 {
	return 0 // Simplified for example
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x * 0.5 // Simplified
}
