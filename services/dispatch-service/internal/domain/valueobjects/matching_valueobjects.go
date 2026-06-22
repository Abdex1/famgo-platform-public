package valueobjects

import (
	"fmt"
	"math"
)

// DriverScore represents all scoring factors for a driver
type DriverScore struct {
	DistanceKm      float64
	Rating          float64      // 1-5 stars
	AcceptanceRate  float64      // 0-1 (0% to 100%)
	Availability    float64      // 0-1
	CompositeScore  float64      // 0-100
}

func NewDriverScore(
	distanceKm float64,
	rating float64,
	acceptedRides int32,
	cancelledRides int32,
) (*DriverScore, error) {
	if distanceKm < 0 {
		return nil, fmt.Errorf("distance cannot be negative: %f", distanceKm)
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be 1-5, got %f", rating)
	}

	totalRides := acceptedRides + cancelledRides
	var acceptanceRate float64
	if totalRides > 0 {
		acceptanceRate = float64(acceptedRides) / float64(totalRides)
	}

	return &DriverScore{
		DistanceKm:     distanceKm,
		Rating:         rating,
		AcceptanceRate: acceptanceRate,
		Availability:   1.0, // Full availability by default
	}, nil
}

// CalculateCompositeScore computes weighted score (0-100)
func (ds *DriverScore) CalculateCompositeScore(
	distanceWeight float64,
	ratingWeight float64,
	acceptanceWeight float64,
	availabilityWeight float64,
) float64 {
	// Normalize factors
	// Distance: inverse relationship (closer = better), max 20km
	distanceFactor := math.Max(0, 1-(ds.DistanceKm/20.0))
	
	// Rating: normalize 1-5 to 0-1
	ratingFactor := (ds.Rating - 1.0) / 4.0
	
	// Acceptance rate: already 0-1
	acceptanceFactor := ds.AcceptanceRate
	
	// Availability: already 0-1
	availabilityFactor := ds.Availability

	// Calculate weighted composite
	composite := (distanceFactor * distanceWeight) +
		(ratingFactor * ratingWeight) +
		(acceptanceFactor * acceptanceWeight) +
		(availabilityFactor * availabilityWeight)

	// Scale to 0-100
	ds.CompositeScore = composite * 100.0
	return ds.CompositeScore
}

// ETA represents estimated time to arrival
type ETA struct {
	DurationMinutes int32
	Distance        float64 // km
	RoutePolyline   string
	Source          string  // GOOGLE_MAPS, CACHED, ESTIMATED
}

func NewETA(durationMinutes int32, distance float64, source string) *ETA {
	return &ETA{
		DurationMinutes: durationMinutes,
		Distance:        distance,
		Source:          source,
	}
}

// Coordinates represents a geographic point
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func NewCoordinates(lat, lng float64) (*Coordinates, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("invalid latitude: %f", lat)
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("invalid longitude: %f", lng)
	}

	return &Coordinates{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

// DistanceTo calculates distance to another coordinate (Haversine)
func (c *Coordinates) DistanceTo(other *Coordinates) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := toRadians(c.Latitude)
	lat2Rad := toRadians(other.Latitude)
	deltaLat := toRadians(other.Latitude - c.Latitude)
	deltaLng := toRadians(other.Longitude - c.Longitude)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	centralAngle := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * centralAngle
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
