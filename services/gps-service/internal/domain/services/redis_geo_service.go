// services/gps-service/internal/domain/services/redis_geo_service.go
// Redis GEO domain service for spatial indexing and queries

package services

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// RedisGeoService provides Redis GEO operations for spatial queries
type RedisGeoService struct {
	// This is a domain service, actual Redis operations handled by infrastructure layer
}

// NewRedisGeoService creates a new Redis GEO service
func NewRedisGeoService() *RedisGeoService {
	return &RedisGeoService{}
}

// GeoLocation represents a location with member ID for Redis storage
type GeoLocation struct {
	MemberID  string
	Latitude  float64
	Longitude float64
}

// GeoRadiusResult represents result of a geo radius query
type GeoRadiusResult struct {
	MemberID string
	Distance float64 // in the specified unit
	Hash     int64
	Latitude float64
	Longitude float64
}

// GeoSearchOptions provides options for geo searches
type GeoSearchOptions struct {
	Longitude      float64
	Latitude       float64
	RadiusKm       float64
	RadiusMiles    float64
	RadiusMeters   float64
	Count          int
	AscendingDist  bool
	WithDistance   bool
	WithHash       bool
	WithCoordinates bool
}

// ValidateGeoLocation validates a geographic location
func (rgs *RedisGeoService) ValidateGeoLocation(geoLoc *GeoLocation) error {
	if geoLoc == nil {
		return fmt.Errorf("geo location is nil")
	}

	if geoLoc.MemberID == "" {
		return fmt.Errorf("member ID cannot be empty")
	}

	if geoLoc.Latitude < -90 || geoLoc.Latitude > 90 {
		return fmt.Errorf("latitude out of valid range: %f", geoLoc.Latitude)
	}

	if geoLoc.Longitude < -180 || geoLoc.Longitude > 180 {
		return fmt.Errorf("longitude out of valid range: %f", geoLoc.Longitude)
	}

	return nil
}

// ConvertDriverLocationToGeoLocation converts driver location to geo location
func (rgs *RedisGeoService) ConvertDriverLocationToGeoLocation(
	driver *entities.DriverLocation,
) *GeoLocation {
	if driver == nil || driver.CurrentLocation == nil {
		return nil
	}

	return &GeoLocation{
		MemberID:  driver.DriverID,
		Latitude:  driver.CurrentLocation.Coordinates.Latitude,
		Longitude: driver.CurrentLocation.Coordinates.Longitude,
	}
}

// ConvertDriverLocationsToGeoLocations converts slice of driver locations
func (rgs *RedisGeoService) ConvertDriverLocationsToGeoLocations(
	drivers []*entities.DriverLocation,
) []*GeoLocation {
	if len(drivers) == 0 {
		return nil
	}

	geoLocations := make([]*GeoLocation, 0, len(drivers))
	for _, driver := range drivers {
		if driver != nil && driver.CurrentLocation != nil {
			geoLocations = append(geoLocations, rgs.ConvertDriverLocationToGeoLocation(driver))
		}
	}

	return geoLocations
}

// BuildGeoSearchOptions builds geo search options with sensible defaults
func (rgs *RedisGeoService) BuildGeoSearchOptions(
	referenceLocation *valueobjects.Geolocation,
	radiusKm float64,
	limit int,
) *GeoSearchOptions {
	if referenceLocation == nil {
		return nil
	}

	return &GeoSearchOptions{
		Latitude:        referenceLocation.Coordinates.Latitude,
		Longitude:       referenceLocation.Coordinates.Longitude,
		RadiusKm:        radiusKm,
		Count:           limit,
		AscendingDist:   true,
		WithDistance:    true,
		WithCoordinates: true,
	}
}

// ValidateGeoSearchOptions validates geo search options
func (rgs *RedisGeoService) ValidateGeoSearchOptions(opts *GeoSearchOptions) error {
	if opts == nil {
		return fmt.Errorf("geo search options is nil")
	}

	if opts.Latitude < -90 || opts.Latitude > 90 {
		return fmt.Errorf("center latitude out of valid range: %f", opts.Latitude)
	}

	if opts.Longitude < -180 || opts.Longitude > 180 {
		return fmt.Errorf("center longitude out of valid range: %f", opts.Longitude)
	}

	if opts.RadiusKm <= 0 && opts.RadiusMiles <= 0 && opts.RadiusMeters <= 0 {
		return fmt.Errorf("at least one radius value must be positive")
	}

	if opts.Count < 0 {
		return fmt.Errorf("count must be non-negative")
	}

	return nil
}

// CalculateSearchRadius calculates optimal search radius based on drivers
func (rgs *RedisGeoService) CalculateSearchRadius(
	expectedDrivers int,
	areaKm2 float64,
) float64 {
	if expectedDrivers <= 0 || areaKm2 <= 0 {
		return 5.0 // Default 5 km
	}

	// Estimate radius from area
	// Area = π * r²
	// r = √(Area / π)
	radiusKm := float64(areaKm2) / 3.14159
	if radiusKm < 0.5 {
		return 0.5
	}
	if radiusKm > 50 {
		return 50
	}

	return radiusKm
}

// PartitionByGeohash partitions locations by geohash prefix
func (rgs *RedisGeoService) PartitionByGeohash(
	locations []*GeoLocation,
	precision int,
) map[string][]*GeoLocation {
	locService := NewLocationService()
	partitions := make(map[string][]*GeoLocation)

	for _, loc := range locations {
		hash, err := locService.CalculateGeohash(loc.Latitude, loc.Longitude, precision)
		if err != nil {
			continue
		}

		prefix := hash[:precision]
		partitions[prefix] = append(partitions[prefix], loc)
	}

	return partitions
}

// MergeGeoResults merges multiple geo search results, deduplicating by member ID
func (rgs *RedisGeoService) MergeGeoResults(
	results ...[]*GeoRadiusResult,
) []*GeoRadiusResult {
	seen := make(map[string]bool)
	var merged []*GeoRadiusResult

	for _, resultSet := range results {
		for _, result := range resultSet {
			if !seen[result.MemberID] {
				merged = append(merged, result)
				seen[result.MemberID] = true
			}
		}
	}

	return merged
}

// SortGeoResultsByDistance sorts geo results by distance
func (rgs *RedisGeoService) SortGeoResultsByDistance(
	results []*GeoRadiusResult,
	ascending bool,
) []*GeoRadiusResult {
	// Bubble sort for simplicity
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			shouldSwap := false
			if ascending {
				shouldSwap = results[j].Distance < results[i].Distance
			} else {
				shouldSwap = results[j].Distance > results[i].Distance
			}

			if shouldSwap {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

// LimitGeoResults limits geo results to specified count
func (rgs *RedisGeoService) LimitGeoResults(
	results []*GeoRadiusResult,
	limit int,
) []*GeoRadiusResult {
	if limit <= 0 || len(results) <= limit {
		return results
	}
	return results[:limit]
}

// CalculateGeoHash calculates 52-bit geohash (Redis native)
func (rgs *RedisGeoService) CalculateGeoHash(
	latitude, longitude float64,
) int64 {
	// Redis uses Geohash with 52-bit precision
	// This is a simplified version - actual Redis implementation is more complex
	lat := latitude
	lng := longitude

	// Normalize to 0-1 range
	latNorm := (lat + 90.0) / 180.0
	lngNorm := (lng + 180.0) / 360.0

	// Interleave bits (simplified - Redis uses more complex algorithm)
	hash := int64(0)
	for i := 0; i < 26; i++ {
		if int(latNorm*float64(1<<uint(26))) & (1 << uint(i)) != 0 {
			hash |= 1 << uint(2*i+1)
		}
		if int(lngNorm*float64(1<<uint(26))) & (1 << uint(i)) != 0 {
			hash |= 1 << uint(2 * i)
		}
	}

	return hash
}

// ValidateGeoArea validates geographic area boundaries
func (rgs *RedisGeoService) ValidateGeoArea(
	centerLat, centerLng, radiusKm float64,
) (bool, error) {
	if radiusKm <= 0 {
		return false, fmt.Errorf("radius must be positive")
	}

	if radiusKm > 6371.0 { // Earth radius
		return false, fmt.Errorf("radius exceeds Earth's radius")
	}

	centerCoords, err := valueobjects.NewCoordinates(centerLat, centerLng)
	if err != nil {
		return false, err
	}

	return centerCoords.IsValid(), nil
}

// IsLocationWithinBounds checks if location is within geographic bounds
func (rgs *RedisGeoService) IsLocationWithinBounds(
	lat, lng float64,
	minLat, maxLat, minLng, maxLng float64,
) bool {
	return lat >= minLat && lat <= maxLat &&
		lng >= minLng && lng <= maxLng
}

// ExpandSearchArea expands search area if not enough results
func (rgs *RedisGeoService) ExpandSearchArea(
	currentRadiusKm float64,
	found int,
	expected int,
	maxRadiusKm float64,
) float64 {
	if found >= expected || currentRadiusKm >= maxRadiusKm {
		return currentRadiusKm
	}

	expansion := float64(expected) / float64(found)
	if expansion > 2.0 {
		expansion = 2.0
	}

	newRadius := currentRadiusKm * expansion
	if newRadius > maxRadiusKm {
		newRadius = maxRadiusKm
	}

	return newRadius
}

// BuildGeohashPrefix builds geohash prefix for indexing
func (rgs *RedisGeoService) BuildGeohashPrefix(
	hash string,
	maxPrefix int,
) string {
	if maxPrefix <= 0 || maxPrefix > len(hash) {
		return hash
	}
	return hash[:maxPrefix]
}
