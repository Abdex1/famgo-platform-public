package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"
)

// DriverCandidate represents ranked driver for dispatch
type DriverCandidate struct {
	DriverID          string
	Latitude          float64
	Longitude         float64
	Rating            float64
	AcceptanceRate    float64
	OnlineDurationSec int64
	ETASeconds        int32
	MatchScore        float64
	Rank              int
}

// DispatchRequest represents ride dispatch job
type DispatchRequest struct {
	ID             string
	RideID         string
	PickupLat      float64
	PickupLng      float64
	DropoffLat     float64
	DropoffLng     float64
	RideType       string
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

// DispatchResult represents dispatch outcome
type DispatchResult struct {
	RideID                string
	AssignedDriverID      string
	Candidates            []DriverCandidate
	DispatchedAt          time.Time
	MatchScore            float64
	EstimatedPickupTime   int32
}

// MatchingAlgorithm implements intelligent driver matching
type MatchingAlgorithm struct {
	// Configuration
	SearchRadiusMeters    float64
	MaxCandidates         int
	MatchTimeoutSeconds   int
	ETA_Weight            float32
	Rating_Weight         float32
	Acceptance_Weight     float32
	OnlineDuration_Weight float32
}

// NewMatchingAlgorithm creates new matching algorithm
func NewMatchingAlgorithm() *MatchingAlgorithm {
	return &MatchingAlgorithm{
		SearchRadiusMeters:    5000, // 5 km
		MaxCandidates:         50,   // Consider top 50 drivers
		MatchTimeoutSeconds:   30,
		ETA_Weight:            0.40, // 40% - most important
		Rating_Weight:         0.30,
		Acceptance_Weight:     0.20,
		OnlineDuration_Weight: 0.10,
	}
}

// FindNearbyDrivers fetches drivers within radius (simulated)
func (m *MatchingAlgorithm) FindNearbyDrivers(ctx context.Context, lat, lng float64) ([]DriverCandidate, error) {
	// In production, query Redis GEO index
	// This is simulated for demonstration
	
	candidates := []DriverCandidate{
		{
			DriverID:       "driver_001",
			Latitude:       lat + 0.001,
			Longitude:      lng + 0.001,
			Rating:         4.8,
			AcceptanceRate: 0.95,
			OnlineDurationSec: 3600,
		},
		{
			DriverID:       "driver_002",
			Latitude:       lat - 0.002,
			Longitude:      lng + 0.002,
			Rating:         4.5,
			AcceptanceRate: 0.90,
			OnlineDurationSec: 7200,
		},
		{
			DriverID:       "driver_003",
			Latitude:       lat + 0.003,
			Longitude:      lng - 0.001,
			Rating:         4.7,
			AcceptanceRate: 0.92,
			OnlineDurationSec: 5400,
		},
	}

	return candidates, nil
}

// CalculateETA calculates estimated time to arrive at pickup
func (m *MatchingAlgorithm) CalculateETA(driverLat, driverLng, pickupLat, pickupLng float64) int32 {
	// Haversine distance calculation
	const R = 6371000 // Earth radius in meters

	lat1 := toRad(driverLat)
	lat2 := toRad(pickupLat)
	deltaLat := toRad(pickupLat - driverLat)
	deltaLng := toRad(pickupLng - driverLng)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c // meters

	// Assume average speed 30 km/h in urban traffic
	// distance (m) / (30 km/h = 8.33 m/s) = time in seconds
	eta := int32(distance / 8.33)
	
	return eta
}

// NormalizeScore normalizes metric to 0-1 range
func (m *MatchingAlgorithm) NormalizeScore(value, min, max float64) float64 {
	if max == min {
		return 0.5
	}
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}
	return normalized
}

// RankDrivers ranks candidates by match score
func (m *MatchingAlgorithm) RankDrivers(candidates []DriverCandidate, pickupLat, pickupLng float64) []DriverCandidate {
	// Calculate scores for all candidates
	etaList := make([]int32, len(candidates))
	ratingList := make([]float64, len(candidates))
	acceptanceList := make([]float64, len(candidates))
	durationList := make([]int64, len(candidates))

	for i, c := range candidates {
		etaList[i] = m.CalculateETA(c.Latitude, c.Longitude, pickupLat, pickupLng)
		ratingList[i] = c.Rating
		acceptanceList[i] = c.AcceptanceRate
		durationList[i] = c.OnlineDurationSec
	}

	// Find min/max for normalization
	maxETA := max(etaList)
	maxRating := max(ratingList)
	maxAcceptance := max(acceptanceList)
	maxDuration := max(durationList)

	// Score each candidate
	for i := range candidates {
		etaScore := 1.0 - m.NormalizeScore(float64(etaList[i]), 0, float64(maxETA))
		ratingScore := m.NormalizeScore(ratingList[i], 0, maxRating)
		acceptanceScore := m.NormalizeScore(acceptanceList[i], 0, maxAcceptance)
		durationScore := m.NormalizeScore(float64(durationList[i]), 0, float64(maxDuration))

		matchScore := etaScore*float64(m.ETA_Weight) +
			ratingScore*float64(m.Rating_Weight) +
			acceptanceScore*float64(m.Acceptance_Weight) +
			durationScore*float64(m.OnlineDuration_Weight)

		candidates[i].ETASeconds = etaList[i]
		candidates[i].MatchScore = matchScore
	}

	// Sort by match score (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MatchScore > candidates[j].MatchScore
	})

	// Assign ranks
	for i := range candidates {
		candidates[i].Rank = i + 1
	}

	return candidates
}

// MatchRide finds best drivers for ride
func (m *MatchingAlgorithm) MatchRide(ctx context.Context, request *DispatchRequest) (*DispatchResult, error) {
	// Step 1: Find nearby drivers
	candidates, err := m.FindNearbyDrivers(ctx, request.PickupLat, request.PickupLng)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby drivers: %w", err)
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no drivers available in area")
	}

	// Step 2: Rank candidates
	rankedCandidates := m.RankDrivers(candidates, request.PickupLat, request.PickupLng)

	// Step 3: Return top candidates (in production, dispatch to top 3)
	topCandidates := rankedCandidates
	if len(topCandidates) > 3 {
		topCandidates = rankedCandidates[:3]
	}

	result := &DispatchResult{
		RideID:              request.RideID,
		AssignedDriverID:    topCandidates[0].DriverID, // Assign best match
		Candidates:          topCandidates,
		DispatchedAt:        time.Now(),
		MatchScore:          topCandidates[0].MatchScore,
		EstimatedPickupTime: topCandidates[0].ETASeconds,
	}

	return result, nil
}

// Helper functions
func toRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func max(list interface{}) float64 {
	switch v := list.(type) {
	case []int32:
		if len(v) == 0 {
			return 0
		}
		m := float64(v[0])
		for _, val := range v[1:] {
			if float64(val) > m {
				m = float64(val)
			}
		}
		return m
	case []float64:
		if len(v) == 0 {
			return 0
		}
		m := v[0]
		for _, val := range v[1:] {
			if val > m {
				m = val
			}
		}
		return m
	case []int64:
		if len(v) == 0 {
			return 0
		}
		m := float64(v[0])
		for _, val := range v[1:] {
			if float64(val) > m {
				m = float64(val)
			}
		}
		return m
	}
	return 0
}
