// services/fraud-service/internal/domain/services/fraud_service.go
package services

import (
	"fmt"
	"math"
	"time"

	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/domain/entities"
)

type FraudService struct {
	highRiskThreshold   float64
	mediumRiskThreshold float64
}

func NewFraudService(highThreshold, mediumThreshold float64) *FraudService {
	return &FraudService{
		highRiskThreshold:   highThreshold,
		mediumRiskThreshold: mediumThreshold,
	}
}

type LocationData struct {
	Latitude    float64
	Longitude   float64
	Timestamp   time.Time
}

func (fs *FraudService) DetectLocationAnomaly(currentLoc LocationData, previousLoc *LocationData) bool {
	if previousLoc == nil {
		return false
	}

	// Calculate distance (simplified Haversine)
	distance := calculateDistance(previousLoc.Latitude, previousLoc.Longitude,
		currentLoc.Latitude, currentLoc.Longitude)

	timeDiff := currentLoc.Timestamp.Sub(previousLoc.Timestamp).Seconds()
	if timeDiff == 0 {
		return false
	}

	// Velocity in km/h
	velocity := (distance / timeDiff) * 3600

	// Impossible velocity: > 900 km/h (speed of commercial airplane)
	if velocity > 900 {
		return true
	}

	return false
}

func (fs *FraudService) DetectVelocityAnomaly(locations []LocationData) bool {
	if len(locations) < 2 {
		return false
	}

	for i := 1; i < len(locations); i++ {
		distance := calculateDistance(locations[i-1].Latitude, locations[i-1].Longitude,
			locations[i].Latitude, locations[i].Longitude)

		timeDiff := locations[i].Timestamp.Sub(locations[i-1].Timestamp).Seconds()
		if timeDiff == 0 {
			continue
		}

		velocity := (distance / timeDiff) * 3600
		if velocity > 150 { // Sustained speed > 150 km/h
			return true
		}
	}

	return false
}

func (fs *FraudService) DetectPaymentAnomaly(currentAmount float64, avgAmount float64, stdDev float64) bool {
	if avgAmount == 0 {
		return false
	}

	// If current amount is more than 3 standard deviations from mean
	zScore := math.Abs((currentAmount - avgAmount) / (stdDev + 0.001))
	return zScore > 3.0
}

func (fs *FraudService) DetectBehaviorAnomaly(rideCount int, avgRideCount int) bool {
	// If user has significantly fewer rides than average (potential fraud)
	if avgRideCount > 0 && rideCount < (avgRideCount / 10) {
		return true
	}
	return false
}

func (fs *FraudService) IsBlacklisted(userID string, blacklist []string) bool {
	for _, id := range blacklist {
		if id == userID {
			return true
		}
	}
	return false
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Asin(math.Sqrt(a))
	return earthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
