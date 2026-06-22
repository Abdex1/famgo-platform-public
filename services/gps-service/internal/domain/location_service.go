package domain

import (
	"math"
)

// LocationService provides pure domain logic for location operations
type LocationService struct {
	// ZERO external dependencies - only domain logic
}

// NewLocationService creates a new location service
func NewLocationService() *LocationService {
	return &LocationService{}
}

// IsWithinGeofence checks if a location is within a geofence (pure logic)
func (s *LocationService) IsWithinGeofence(location DriverLocation, geofence Geofence) bool {
	distance := s.CalculateDistance(location.Latitude, location.Longitude, geofence.Latitude, geofence.Longitude)
	return distance <= float64(geofence.Radius)
}

// CalculateDistance calculates distance between two coordinates in meters (Haversine formula)
// Pure domain logic - no external dependencies
func (s *LocationService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6371000.0 // Earth's radius in meters

	// Convert to radians
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

// CalculateDeviation calculates deviation between expected and actual location (pure logic)
func (s *LocationService) CalculateDeviation(expected, actual DriverLocation) float64 {
	return s.CalculateDistance(expected.Latitude, expected.Longitude, actual.Latitude, actual.Longitude)
}

// IsSignificantDeviation checks if deviation exceeds threshold (pure logic)
func (s *LocationService) IsSignificantDeviation(expected, actual DriverLocation, thresholdMeters float64) bool {
	deviation := s.CalculateDeviation(expected, actual)
	return deviation > thresholdMeters
}

// CalculateETA calculates estimated time of arrival based on distance and average speed (pure logic)
func (s *LocationService) CalculateETA(distanceMeters float64, averageSpeedKmH float64) int {
	if averageSpeedKmH <= 0 {
		return 0
	}
	distanceKm := distanceMeters / 1000.0
	hours := distanceKm / averageSpeedKmH
	minutes := int(hours * 60)
	return minutes
}

// Helper function to convert degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
