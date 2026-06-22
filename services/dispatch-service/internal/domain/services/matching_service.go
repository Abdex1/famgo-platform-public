// services/dispatch-service/internal/domain/services/matching_service.go
// MatchingService with multi-factor driver scoring algorithm

package services

import (
	"fmt"
	"math"
	"sort"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/valueobjects"
)

// DriverInfo holds driver data for matching
type DriverInfo struct {
	DriverID         string
	Latitude         float64
	Longitude        float64
	IsOnline         bool
	AcceptanceRate   float64 // 0-100 percentage
	Rating           float64 // 0-5 stars
	Distance         float64 // kilometers from pickup
	ETA              float64 // minutes
}

// MatchingService provides driver matching with multi-factor scoring
type MatchingService struct {
	proximityWeight       float64
	acceptanceRateWeight  float64
	ratingWeight          float64
	availabilityWeight    float64
	minAcceptanceRate     float64
	minRating             float64
	maxDistanceKm         float64
}

// NewMatchingService creates a new matching service
func NewMatchingService(
	proximityWeight, acceptanceRateWeight, ratingWeight, availabilityWeight,
	minAcceptanceRate, minRating, maxDistanceKm float64,
) *MatchingService {
	return &MatchingService{
		proximityWeight:      proximityWeight,
		acceptanceRateWeight: acceptanceRateWeight,
		ratingWeight:         ratingWeight,
		availabilityWeight:   availabilityWeight,
		minAcceptanceRate:    minAcceptanceRate,
		minRating:            minRating,
		maxDistanceKm:        maxDistanceKm,
	}
}

// MatchDrivers scores and ranks drivers for matching
func (ms *MatchingService) MatchDrivers(
	pickupLat, pickupLng float64,
	drivers []*DriverInfo,
	limit int,
) ([]*valueobjects.MatchScore, error) {
	if pickupLat < -90 || pickupLat > 90 {
		return nil, fmt.Errorf("invalid pickup latitude: %f", pickupLat)
	}
	if pickupLng < -180 || pickupLng > 180 {
		return nil, fmt.Errorf("invalid pickup longitude: %f", pickupLng)
	}

	if len(drivers) == 0 {
		return []*valueobjects.MatchScore{}, nil
	}

	var scores []*valueobjects.MatchScore

	for _, driver := range drivers {
		// Filter out drivers that don't meet minimum criteria
		if !driver.IsOnline {
			continue
		}
		if driver.AcceptanceRate < ms.minAcceptanceRate {
			continue
		}
		if driver.Rating < ms.minRating {
			continue
		}
		if driver.Distance > ms.maxDistanceKm {
			continue
		}

		// Calculate individual scores
		proximityScore := ms.calculateProximityScore(driver.Distance)
		acceptanceScore := driver.AcceptanceRate // Already 0-100
		ratingScore := (driver.Rating / 5.0) * 100 // Normalize to 0-100
		availabilityScore := 100.0 // Online drivers get full score

		// Create match score
		matchScore, err := valueobjects.NewMatchScore(
			driver.DriverID,
			proximityScore,
			acceptanceScore,
			ratingScore,
			availabilityScore,
			ms.proximityWeight,
			ms.acceptanceRateWeight,
			ms.ratingWeight,
			ms.availabilityWeight,
			driver.Distance,
			driver.ETA,
		)
		if err != nil {
			continue
		}

		scores = append(scores, matchScore)
	}

	// Sort by total score (highest first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})

	// Limit results
	if limit > 0 && len(scores) > limit {
		scores = scores[:limit]
	}

	return scores, nil
}

// calculateProximityScore calculates proximity score based on distance
// Closer drivers get higher scores (0-100)
func (ms *MatchingService) calculateProximityScore(distanceKm float64) float64 {
	if distanceKm <= 0 {
		return 100.0
	}

	// Linear decay: at maxDistance, score is 0
	if distanceKm >= ms.maxDistanceKm {
		return 0.0
	}

	score := (1 - (distanceKm / ms.maxDistanceKm)) * 100
	return math.Max(0, math.Min(100, score))
}

// CalculateOptimalSearchRadius calculates optimal search radius based on available drivers
func (ms *MatchingService) CalculateOptimalSearchRadius(
	availableDriverCount int,
	targetMatchCount int,
	currentRadius float64,
	maxRadius float64,
) float64 {
	if availableDriverCount >= targetMatchCount {
		return currentRadius
	}

	if currentRadius >= maxRadius {
		return maxRadius
	}

	// Expand radius proportionally
	expansionFactor := float64(targetMatchCount) / float64(availableDriverCount)
	if expansionFactor > 2.0 {
		expansionFactor = 2.0 // Cap at 2x expansion per iteration
	}

	newRadius := currentRadius * expansionFactor
	if newRadius > maxRadius {
		newRadius = maxRadius
	}

	return newRadius
}

// ValidateDriversForMatching validates driver eligibility
func (ms *MatchingService) ValidateDriversForMatching(driver *DriverInfo) (bool, string) {
	if !driver.IsOnline {
		return false, "driver offline"
	}

	if driver.AcceptanceRate < ms.minAcceptanceRate {
		return false, fmt.Sprintf("acceptance rate %.1f%% below minimum %.1f%%",
			driver.AcceptanceRate, ms.minAcceptanceRate)
	}

	if driver.Rating < ms.minRating {
		return false, fmt.Sprintf("rating %.1f below minimum %.1f",
			driver.Rating, ms.minRating)
	}

	if driver.Distance > ms.maxDistanceKm {
		return false, fmt.Sprintf("distance %.2f km exceeds max %.2f km",
			driver.Distance, ms.maxDistanceKm)
	}

	return true, ""
}

// ScoreComparison compares two match scores
type ScoreComparison struct {
	DriverA              string
	DriverB              string
	Winner               string
	WinnerScore          float64
	DifferenceFactor     float64
	ReasonForWin         string
}

// CompareMatches compares two drivers for matching
func (ms *MatchingService) CompareMatches(
	scoreA, scoreB *valueobjects.MatchScore,
) *ScoreComparison {
	comparison := &ScoreComparison{
		DriverA:         scoreA.DriverID,
		DriverB:         scoreB.DriverID,
		DifferenceFactor: scoreA.TotalScore / math.Max(1, scoreB.TotalScore),
	}

	if scoreA.TotalScore > scoreB.TotalScore {
		comparison.Winner = scoreA.DriverID
		comparison.WinnerScore = scoreA.TotalScore

		// Determine reason for win
		if scoreA.ProximityScore > scoreB.ProximityScore+5 {
			comparison.ReasonForWin = "better proximity"
		} else if scoreA.AcceptanceRateScore > scoreB.AcceptanceRateScore+5 {
			comparison.ReasonForWin = "better acceptance rate"
		} else if scoreA.RatingScore > scoreB.RatingScore+5 {
			comparison.ReasonForWin = "better rating"
		} else {
			comparison.ReasonForWin = "overall score"
		}
	} else {
		comparison.Winner = scoreB.DriverID
		comparison.WinnerScore = scoreB.TotalScore

		if scoreB.ProximityScore > scoreA.ProximityScore+5 {
			comparison.ReasonForWin = "better proximity"
		} else if scoreB.AcceptanceRateScore > scoreA.AcceptanceRateScore+5 {
			comparison.ReasonForWin = "better acceptance rate"
		} else if scoreB.RatingScore > scoreA.RatingScore+5 {
			comparison.ReasonForWin = "better rating"
		} else {
			comparison.ReasonForWin = "overall score"
		}
	}

	return comparison
}

// GetScoreBreakdown returns a human-readable breakdown of scores
func (ms *MatchingService) GetScoreBreakdown(score *valueobjects.MatchScore) map[string]interface{} {
	return map[string]interface{}{
		"driver_id":             score.DriverID,
		"total_score":           fmt.Sprintf("%.2f", score.TotalScore),
		"proximity_score":       fmt.Sprintf("%.2f (weight: %.0f%%)", score.ProximityScore, score.ProximityWeight*100),
		"acceptance_score":      fmt.Sprintf("%.2f (weight: %.0f%%)", score.AcceptanceRateScore, score.AcceptanceWeight*100),
		"rating_score":          fmt.Sprintf("%.2f (weight: %.0f%%)", score.RatingScore, score.RatingWeight*100),
		"availability_score":    fmt.Sprintf("%.2f (weight: %.0f%%)", score.AvailabilityScore, score.AvailabilityWeight*100),
		"distance_km":           fmt.Sprintf("%.2f", score.Distance),
		"eta_minutes":           fmt.Sprintf("%.0f", score.ETA),
	}
}
