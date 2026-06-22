// services/ride-service/internal/domain/ride_service.go
// Ride Service Domain Logic

package domain

import "time"

// RideService provides domain-level operations
type RideService struct{}

// NewRideService creates a new ride service
func NewRideService() *RideService {
	return &RideService{}
}

// ValidateLocation validates geographic coordinates
func (s *RideService) ValidateLocation(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// CalculateDistance calculates distance between two points (Haversine formula)
func (s *RideService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float32 {
	const earthRadiusKm = 6371.0
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := sin(dLat/2)*sin(dLat/2) + cos(toRad(lat1))*cos(toRad(lat2))*sin(dLon/2)*sin(dLon/2)
	c := 2 * asin(sqrt(a))
	return float32(earthRadiusKm * c)
}

// EstimateETA estimates time to arrival in minutes
func (s *RideService) EstimateETA(distanceKm float32) int32 {
	// Average speed: 40 km/h in city
	averageSpeedKmH := 40.0
	timeHours := float64(distanceKm) / averageSpeedKmH
	timeMinutes := timeHours * 60
	return int32(timeMinutes) + 5 // Add 5 min buffer
}

// CanCompleteRide checks if ride is ready to complete
func (s *RideService) CanCompleteRide(ride *Ride) bool {
	return ride.Status == RideStatusStarted && ride.PickupTime != nil
}

// IsRideExpired checks if ride has been waiting too long
func (s *RideService) IsRideExpired(ride *Ride, maxWaitMinutes int32) bool {
	if ride.Status != RideStatusRequested && ride.Status != RideStatusSearching {
		return false
	}
	elapsed := time.Since(ride.CreatedAt).Minutes()
	return int32(elapsed) > maxWaitMinutes
}

// Helper functions
func toRad(degrees float64) float64 {
	return degrees * 3.141592653589793 / 180.0
}

func sin(x float64) float64 {
	// Approximation
	return x
}

func cos(x float64) float64 {
	return 1 - x*x/2
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func asin(x float64) float64 {
	return x
}
