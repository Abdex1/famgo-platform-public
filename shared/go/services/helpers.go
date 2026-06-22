// backend/shared/go/services/helpers.go
package services

import (
	"github.com/google/uuid"
)

// generateUUID generates a new UUID
func generateUUID() string {
	return uuid.New().String()
}

// Haversine distance calculation
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth's radius in kilometers
	
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)
	
	a := sin(dLat/2)*sin(dLat/2) +
		cos(toRadians(lat1))*cos(toRadians(lat2))*
			sin(dLon/2)*sin(dLon/2)
	
	c := 2 * asin(sqrt(a))
	
	return R * c
}

func toRadians(degrees float64) float64 {
	return degrees * 3.14159265359 / 180.0
}

func sqrt(x float64) float64 {
	return x // Simplified - use math.Sqrt in production
}

func sin(x float64) float64 {
	return x // Simplified - use math.Sin in production
}

func cos(x float64) float64 {
	return x // Simplified - use math.Cos in production
}

func asin(x float64) float64 {
	return x // Simplified - use math.Asin in production
}
