package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/domain/entities"
)

// PoolingEngine handles ride pooling logic
type PoolingEngine struct {
	// Weights for compatibility scoring
	RouteOverlapWeight    float64
	ProfitabilityWeight   float64
	ETASimilarityWeight   float64
	PickupProximityWeight float64

	// Constraints
	MaxDetourMinutes   int
	MaxWaitMinutes     int
	MinRouteOverlap    float64
	MaxPoolSize        int
	MinPoolSize        int
}

// NewPoolingEngine creates new pooling engine
func NewPoolingEngine() *PoolingEngine {
	return &PoolingEngine{
		RouteOverlapWeight:    0.4,
		ProfitabilityWeight:   0.3,
		ETASimilarityWeight:   0.2,
		PickupProximityWeight: 0.1,
		MaxDetourMinutes:      10,
		MaxWaitMinutes:        5,
		MinRouteOverlap:       0.70,
		MaxPoolSize:           3,
		MinPoolSize:           2,
	}
}

// CalculateRouteOverlap calculates overlap between two routes (0-1)
func (e *PoolingEngine) CalculateRouteOverlap(
	pickup1Lat, pickup1Lng, dropoff1Lat, dropoff1Lng,
	pickup2Lat, pickup2Lng, dropoff2Lat, dropoff2Lng float64,
) float64 {
	// Calculate bounding box overlap
	// Route 1 bounding box
	minLat1 := math.Min(pickup1Lat, dropoff1Lat)
	maxLat1 := math.Max(pickup1Lat, dropoff1Lat)
	minLng1 := math.Min(pickup1Lng, dropoff1Lng)
	maxLng1 := math.Max(pickup1Lng, dropoff1Lng)

	// Route 2 bounding box
	minLat2 := math.Min(pickup2Lat, dropoff2Lat)
	maxLat2 := math.Max(pickup2Lat, dropoff2Lat)
	minLng2 := math.Min(pickup2Lng, dropoff2Lng)
	maxLng2 := math.Max(pickup2Lng, dropoff2Lng)

	// Calculate overlap area
	overlapLat := math.Min(maxLat1, maxLat2) - math.Max(minLat1, minLat2)
	overlapLng := math.Min(maxLng1, maxLng2) - math.Max(minLng1, minLng2)

	if overlapLat < 0 || overlapLng < 0 {
		return 0.0 // No overlap
	}

	overlapArea := overlapLat * overlapLng

	// Total area covered by both routes
	totalArea := (maxLat1-minLat1)*(maxLng1-minLng1) + (maxLat2-minLat2)*(maxLng2-minLng2) - overlapArea

	if totalArea == 0 {
		return 0.0
	}

	return overlapArea / totalArea
}

// CalculateDetour estimates extra travel time for pick up
func (e *PoolingEngine) CalculateDetour(
	pickupLat, pickupLng,
	route1PickupLat, route1PickupLng, route1DropoffLat, route1DropoffLng,
	route2PickupLat, route2PickupLng, route2DropoffLat, route2DropoffLng float64,
) int {
	// Simplified: use straight-line distance differences
	// In production: use Google Maps API

	// Distance from route 1 dropoff to route 2 pickup
	directDistance := e.haversineDistance(
		route1DropoffLat, route1DropoffLng,
		route2PickupLat, route2PickupLng,
	)

	// Detour is extra distance if we go via the new pickup
	detourDistance := directDistance - 5000 // Assume 5km baseline

	// Convert to minutes (assume 30 km/h average)
	detourMinutes := int(detourDistance / 500) // 30 km/h = 500 m/min

	if detourMinutes < 0 {
		detourMinutes = 0
	}

	return detourMinutes
}

// CheckPoolConstraints validates pool against constraints
func (e *PoolingEngine) CheckPoolConstraints(
	ride1, ride2 *entities.PoolRequest,
) (valid bool, reason string) {
	// Check detour time
	detour := e.CalculateDetour(
		ride1.PickupLat, ride1.PickupLng,
		ride1.PickupLat, ride1.PickupLng, ride1.DropoffLat, ride1.DropoffLng,
		ride2.PickupLat, ride2.PickupLng, ride2.DropoffLat, ride2.DropoffLng,
	)

	if detour > e.MaxDetourMinutes {
		return false, fmt.Sprintf("detour exceeds limit: %d > %d min", detour, e.MaxDetourMinutes)
	}

	// Check wait time for existing passenger
	waitTime := int(ride2.CreatedAt.Sub(ride1.CreatedAt).Minutes())
	if waitTime > e.MaxWaitMinutes {
		return false, fmt.Sprintf("wait exceeds limit: %d > %d min", waitTime, e.MaxWaitMinutes)
	}

	// Check route overlap
	overlap := e.CalculateRouteOverlap(
		ride1.PickupLat, ride1.PickupLng, ride1.DropoffLat, ride1.DropoffLng,
		ride2.PickupLat, ride2.PickupLng, ride2.DropoffLat, ride2.DropoffLng,
	)

	if overlap < e.MinRouteOverlap {
		return false, fmt.Sprintf("route overlap too low: %.2f < %.2f", overlap, e.MinRouteOverlap)
	}

	return true, ""
}

// CalculateCompatibility scores two rides for pooling
func (e *PoolingEngine) CalculateCompatibility(
	ride1, ride2 *entities.PoolRequest,
) *entities.PoolCompatibility {
	compat := &entities.PoolCompatibility{}

	// 1. Route overlap (weight: 0.4)
	compat.RouteOverlapScore = e.CalculateRouteOverlap(
		ride1.PickupLat, ride1.PickupLng, ride1.DropoffLat, ride1.DropoffLng,
		ride2.PickupLat, ride2.PickupLng, ride2.DropoffLat, ride2.DropoffLng,
	)

	// 2. Profitability (weight: 0.3) - Higher fares are more profitable
	// Assume pooled discount is 20%, so profitability = fare retention
	avgFare := (ride1.EstimatedFare + ride2.EstimatedFare) / 2
	if avgFare > 100 {
		compat.ProfitabilityScore = 1.0 // Premium rides
	} else if avgFare > 50 {
		compat.ProfitabilityScore = 0.8
	} else {
		compat.ProfitabilityScore = 0.6 // Cheaper rides less profitable
	}

	// 3. ETA similarity (weight: 0.2) - Similar duration means compatible
	durationDiff := math.Abs(float64(ride1.EstimatedDuration - ride2.EstimatedDuration))
	avgDuration := float64(ride1.EstimatedDuration+ride2.EstimatedDuration) / 2
	etaDiffRatio := durationDiff / avgDuration
	if etaDiffRatio > 0.5 {
		compat.ETASimilarityScore = 0.5 // Very different
	} else if etaDiffRatio > 0.2 {
		compat.ETASimilarityScore = 0.7
	} else {
		compat.ETASimilarityScore = 0.9 // Very similar
	}

	// 4. Pickup proximity (weight: 0.1) - Closer pickups are better
	pickupDistance := e.haversineDistance(
		ride1.PickupLat, ride1.PickupLng,
		ride2.PickupLat, ride2.PickupLng,
	)
	if pickupDistance > 2000 {
		compat.PickupProximityScore = 0.3 // Far apart
	} else if pickupDistance > 1000 {
		compat.PickupProximityScore = 0.6
	} else {
		compat.PickupProximityScore = 1.0 // Very close
	}

	// Calculate final weighted score
	compat.FinalScore = (compat.RouteOverlapScore * e.RouteOverlapWeight) +
		(compat.ProfitabilityScore * e.ProfitabilityWeight) +
		(compat.ETASimilarityScore * e.ETASimilarityWeight) +
		(compat.PickupProximityScore * e.PickupProximityWeight)

	// Viable if score > 0.5 AND route overlap >= min threshold
	compat.IsViable = compat.FinalScore > 0.5 && compat.RouteOverlapScore >= e.MinRouteOverlap

	return compat
}

// FindPoolCandidates finds compatible rides for pooling
func (e *PoolingEngine) FindPoolCandidates(
	targetRide *entities.PoolRequest,
	availableRides []*entities.PoolRequest,
) []entities.PoolCandidate {
	var candidates []entities.PoolCandidate

	for _, candidate := range availableRides {
		// Skip same ride
		if candidate.RideID == targetRide.RideID {
			continue
		}

		// Check basic constraints
		valid, _ := e.CheckPoolConstraints(targetRide, candidate)
		if !valid {
			continue
		}

		// Calculate compatibility
		compat := e.CalculateCompatibility(targetRide, candidate)
		if !compat.IsViable {
			continue
		}

		// Calculate savings (20-30% typical for pooled rides)
		savingsPercentage := 0.25 // 25% average discount
		detour := e.CalculateDetour(
			targetRide.PickupLat, targetRide.PickupLng,
			targetRide.PickupLat, targetRide.PickupLng, targetRide.DropoffLat, targetRide.DropoffLng,
			candidate.PickupLat, candidate.PickupLng, candidate.DropoffLat, candidate.DropoffLng,
		)

		candidates = append(candidates, entities.PoolCandidate{
			RideID:             candidate.RideID,
			PickupLat:          candidate.PickupLat,
			PickupLng:          candidate.PickupLng,
			DropoffLat:         candidate.DropoffLat,
			DropoffLng:         candidate.DropoffLng,
			PickupAddress:      candidate.PickupAddress,
			DropoffAddress:     candidate.DropoffAddress,
			RouteOverlap:       compat.RouteOverlapScore,
			DetourMinutes:      detour,
			WaitMinutes:        int(candidate.CreatedAt.Sub(targetRide.CreatedAt).Minutes()),
			CompatibilityScore: compat.FinalScore,
			SavingsPercentage:  savingsPercentage,
		})
	}

	// Sort by compatibility score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].CompatibilityScore > candidates[j].CompatibilityScore
	})

	// Assign ranks
	for i := range candidates {
		candidates[i].Rank = i + 1
	}

	return candidates
}

// OptimizeRoute calculates optimal pickup/dropoff sequence
func (e *PoolingEngine) OptimizeRoute(
	rides []*entities.PoolRequest,
) []string {
	// Simple algorithm: nearest neighbor
	// In production: use traveling salesman problem solver

	if len(rides) == 0 {
		return []string{}
	}

	sequence := make([]string, 0, len(rides))
	sequence = append(sequence, rides[0].RideID)

	for len(sequence) < len(rides) {
		lastRide := rides[0]
		for _, r := range rides {
			if r.RideID == sequence[len(sequence)-1] {
				lastRide = r
				break
			}
		}

		// Find nearest unvisited
		minDist := math.MaxFloat64
		nextIdx := -1
		for i, r := range rides {
			// Check if already in sequence
			inSeq := false
			for _, s := range sequence {
				if s == r.RideID {
					inSeq = true
					break
				}
			}
			if inSeq {
				continue
			}

			dist := e.haversineDistance(
				lastRide.DropoffLat, lastRide.DropoffLng,
				r.PickupLat, r.PickupLng,
			)

			if dist < minDist {
				minDist = dist
				nextIdx = i
			}
		}

		if nextIdx >= 0 {
			sequence = append(sequence, rides[nextIdx].RideID)
		}
	}

	return sequence
}

// Helper: Haversine distance
func (e *PoolingEngine) haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // meters
	const toRad = math.Pi / 180

	dLat := (lat2 - lat1) * toRad
	dLng := (lng2 - lng1) * toRad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
