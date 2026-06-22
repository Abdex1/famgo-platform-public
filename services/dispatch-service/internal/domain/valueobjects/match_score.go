// services/dispatch-service/internal/domain/valueobjects/match_score.go
// MatchScore value object for driver scoring

package valueobjects

import (
	"fmt"
)

// MatchScore holds the detailed scoring breakdown for a driver match
type MatchScore struct {
	DriverID              string
	ProximityScore        float64 // 0-100, distance-based
	AcceptanceRateScore   float64 // 0-100, acceptance rate-based
	RatingScore           float64 // 0-100, rating-based
	AvailabilityScore     float64 // 0-100, online status-based
	TotalScore            float64 // 0-100, weighted sum
	Distance              float64 // kilometers
	ETA                   float64 // minutes
	ProximityWeight       float64
	AcceptanceWeight      float64
	RatingWeight          float64
	AvailabilityWeight    float64
}

// NewMatchScore creates a new match score
func NewMatchScore(
	driverID string,
	proximityScore, acceptanceScore, ratingScore, availabilityScore float64,
	proximityWeight, acceptanceWeight, ratingWeight, availabilityWeight float64,
	distance, eta float64,
) (*MatchScore, error) {
	if driverID == "" {
		return nil, fmt.Errorf("driver ID cannot be empty")
	}

	// Clamp all scores to 0-100 range
	ps := clamp(proximityScore, 0, 100)
	as := clamp(acceptanceScore, 0, 100)
	rs := clamp(ratingScore, 0, 100)
	avs := clamp(availabilityScore, 0, 100)

	// Calculate weighted total
	totalWeight := proximityWeight + acceptanceWeight + ratingWeight + availabilityWeight
	if totalWeight <= 0 {
		return nil, fmt.Errorf("total weight must be positive")
	}

	totalScore := (ps*proximityWeight +
		as*acceptanceWeight +
		rs*ratingWeight +
		avs*availabilityWeight) / totalWeight

	return &MatchScore{
		DriverID:            driverID,
		ProximityScore:      ps,
		AcceptanceRateScore: as,
		RatingScore:         rs,
		AvailabilityScore:   avs,
		TotalScore:          totalScore,
		Distance:            distance,
		ETA:                 eta,
		ProximityWeight:     proximityWeight,
		AcceptanceWeight:    acceptanceWeight,
		RatingWeight:        ratingWeight,
		AvailabilityWeight:  availabilityWeight,
	}, nil
}

// IsValid checks if score is valid
func (ms *MatchScore) IsValid() bool {
	return ms.DriverID != "" &&
		ms.TotalScore >= 0 &&
		ms.TotalScore <= 100 &&
		ms.Distance >= 0 &&
		ms.ETA >= 0
}

// String returns string representation
func (ms *MatchScore) String() string {
	return fmt.Sprintf(
		"Driver: %s, Total: %.2f, Proximity: %.2f, Acceptance: %.2f, Rating: %.2f, Availability: %.2f, Distance: %.2fkm, ETA: %.0fmin",
		ms.DriverID, ms.TotalScore, ms.ProximityScore, ms.AcceptanceRateScore,
		ms.RatingScore, ms.AvailabilityScore, ms.Distance, ms.ETA,
	)
}

// clamp clamps a value between min and max
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
