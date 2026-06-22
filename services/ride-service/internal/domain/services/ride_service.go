// services/ride-service/internal/domain/services/ride_service.go
// Ride domain service for business logic

package services

import (
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/entities"
)

// RideService provides ride-related business logic
type RideService struct{}

// NewRideService creates a new ride service
func NewRideService() *RideService {
	return &RideService{}
}

// CalculateFare calculates ride fare based on distance and duration
type FareCalculation struct {
	BaseFare          float64
	DistanceFare      float64
	TimeFare          float64
	SurgePricing      float64
	Discount          float64
	TotalFare         float64
}

// CalculateFare calculates fare for a ride
func (rs *RideService) CalculateFare(
	distance float64,
	duration time.Duration,
	surgeMultiplier float64,
	discount float64,
) *FareCalculation {
	const (
		baseFare           = 50.0      // Birr
		distanceRate       = 8.0       // Birr per km
		timeRate           = 0.5       // Birr per minute
		minimumFare        = 80.0      // Birr
	)

	// Base fare
	base := baseFare

	// Distance fare
	distanceFare := distance * distanceRate

	// Time fare
	timeFare := duration.Minutes() * timeRate

	// Subtotal
	subtotal := base + distanceFare + timeFare

	// Surge pricing
	surgeFare := subtotal * (surgeMultiplier - 1.0)
	if surgeFare < 0 {
		surgeFare = 0
	}

	// Subtotal with surge
	withSurge := subtotal + surgeFare

	// Apply discount
	discountAmount := withSurge * discount
	total := withSurge - discountAmount

	// Apply minimum fare
	if total < minimumFare {
		total = minimumFare
	}

	return &FareCalculation{
		BaseFare:     base,
		DistanceFare: distanceFare,
		TimeFare:     timeFare,
		SurgePricing: surgeFare,
		Discount:     discountAmount,
		TotalFare:    total,
	}
}

// ValidateRideRequest validates ride request
func (rs *RideService) ValidateRideRequest(
	pickupLat, pickupLng, dropoffLat, dropoffLng float64,
	minDistance, maxDistance float64,
) (bool, string) {
	// Validate coordinates
	if pickupLat < -90 || pickupLat > 90 {
		return false, "invalid pickup latitude"
	}
	if pickupLng < -180 || pickupLng > 180 {
		return false, "invalid pickup longitude"
	}
	if dropoffLat < -90 || dropoffLat > 90 {
		return false, "invalid dropoff latitude"
	}
	if dropoffLng < -180 || dropoffLng > 180 {
		return false, "invalid dropoff longitude"
	}

	// Calculate distance
	distance := haversine(pickupLat, pickupLng, dropoffLat, dropoffLng)

	// Validate distance
	if distance < minDistance {
		return false, fmt.Sprintf("distance too short (%.2f km, minimum %.2f km)", distance, minDistance)
	}
	if distance > maxDistance {
		return false, fmt.Sprintf("distance too long (%.2f km, maximum %.2f km)", distance, maxDistance)
	}

	return true, ""
}

// CanTransitionToStatus checks if status transition is valid
func (rs *RideService) CanTransitionToStatus(currentStatus, newStatus entities.RideStatus) bool {
	transitions := map[entities.RideStatus][]entities.RideStatus{
		entities.StatusRequested: {entities.StatusAccepted, entities.StatusCancelled},
		entities.StatusAccepted:  {entities.StatusPickedUp, entities.StatusCancelled, entities.StatusNoShow},
		entities.StatusPickedUp:  {entities.StatusOngoing, entities.StatusCancelled, entities.StatusNoShow},
		entities.StatusOngoing:   {entities.StatusCompleted},
		entities.StatusCompleted: {},
		entities.StatusCancelled: {},
		entities.StatusNoShow:    {},
	}

	allowedTransitions, exists := transitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// CalculateWaitTime estimates wait time based on driver location and pickup
func (rs *RideService) CalculateWaitTime(
	driverDistance float64,
	avgSpeedKmH float64,
) time.Duration {
	if driverDistance <= 0 || avgSpeedKmH <= 0 {
		return time.Duration(0)
	}

	// distance / speed = time in hours
	hours := driverDistance / avgSpeedKmH
	minutes := hours * 60

	return time.Duration(minutes) * time.Minute
}

// RideMetrics holds ride performance metrics
type RideMetrics struct {
	RideID       string
	Status       string
	Distance     float64
	Duration     time.Duration
	AverageSpeed float64
	Rating       int
	Fare         float64
}

// CalculateMetrics calculates ride metrics
func (rs *RideService) CalculateMetrics(ride *entities.Ride) *RideMetrics {
	distance := ride.GetDistance()
	duration := ride.GetDuration()

	var avgSpeed float64
	if duration > 0 {
		avgSpeed = distance / (duration / 60) // distance / hours
	}

	var rating int
	if ride.DriverRating != nil {
		rating = *ride.DriverRating
	}

	var fare float64
	if ride.ActualFare != nil {
		fare = *ride.ActualFare
	} else {
		fare = ride.EstimatedFare
	}

	return &RideMetrics{
		RideID:       ride.ID,
		Status:       string(ride.Status),
		Distance:     distance,
		Duration:     time.Duration(duration) * time.Minute,
		AverageSpeed: avgSpeed,
		Rating:       rating,
		Fare:         fare,
	}
}

// Helper function: Haversine formula for distance
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := toRadians(lat1)
	lat2Rad := toRadians(lat2)
	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := sin2(dLat/2) + cos(lat1Rad)*cos(lat2Rad)*sin2(dLng/2)
	c := 2 * asin(sqrt(a))

	return earthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return degrees * 3.14159265359 / 180.0
}

func sin2(x float64) float64 {
	s := sin(x)
	return s * s
}

func sin(x float64) float64 {
	// Simplified sine
	return x - (x*x*x)/6 + (x*x*x*x*x)/120
}

func cos(x float64) float64 {
	return 1 - (x*x)/2 + (x*x*x*x)/24
}

func asin(x float64) float64 {
	return x + (x*x*x)/6 + (3*x*x*x*x*x)/40
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
