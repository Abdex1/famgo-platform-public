package ports

import "context"

// DriverCandidate is the normalized driver snapshot used by dispatch matching.
type DriverCandidate struct {
	DriverID       string
	Latitude       float64
	Longitude      float64
	DistanceKm     float64
	ETAMinutes     float64
	IsOnline       bool
	AcceptanceRate float64
	Rating         float64
}

// DriverDiscovery finds eligible drivers near a pickup location.
type DriverDiscovery interface {
	FindDriversWithinRadius(
		ctx context.Context,
		latitude, longitude, radiusKm float64,
		limit int,
	) ([]DriverCandidate, error)
}
