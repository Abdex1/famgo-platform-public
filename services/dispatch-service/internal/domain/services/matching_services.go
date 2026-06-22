package services

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/valueobjects"
)

// MatchingAlgorithmService implements the driver-rider matching logic
type MatchingAlgorithmService struct {
	distanceWeight      float64
	ratingWeight        float64
	acceptanceWeight    float64
	availabilityWeight  float64
}

func NewMatchingAlgorithmService(
	distWeight, ratingWeight, acceptanceWeight, availWeight float64,
) *MatchingAlgorithmService {
	return &MatchingAlgorithmService{
		distanceWeight:     distWeight,
		ratingWeight:       ratingWeight,
		acceptanceWeight:   acceptanceWeight,
		availabilityWeight: availWeight,
	}
}

// RankDrivers scores and ranks drivers by suitability
func (mas *MatchingAlgorithmService) RankDrivers(
	ctx context.Context,
	drivers []*entities.DriverMatch,
	pickupLat, pickupLng float64,
	acceptanceThreshold float64,
) ([]*entities.DriverMatch, error) {
	if len(drivers) == 0 {
		return nil, fmt.Errorf("no drivers available")
	}

	pickup, _ := valueobjects.NewCoordinates(pickupLat, pickupLng)

	// Score each driver
	rankedDrivers := make([]*entities.DriverMatch, 0, len(drivers))
	for _, driver := range drivers {
		// Skip drivers below acceptance threshold
		acceptanceRate := 0.0
		if driver.AcceptedRides+driver.CancelledRides > 0 {
			acceptanceRate = float64(driver.AcceptedRides) / float64(driver.AcceptedRides+driver.CancelledRides)
		}
		if acceptanceRate < acceptanceThreshold {
			continue
		}

		// Calculate composite score
		driverCoord, _ := valueobjects.NewCoordinates(driver.Latitude, driver.Longitude)
		distance := pickup.DistanceTo(driverCoord)

		dScore, _ := valueobjects.NewDriverScore(
			distance,
			driver.Rating,
			driver.AcceptedRides,
			driver.CancelledRides,
		)

		score := dScore.CalculateCompositeScore(
			mas.distanceWeight,
			mas.ratingWeight,
			mas.acceptanceWeight,
			mas.availabilityWeight,
		)

		driver.Score = score
		driver.Distance = distance
		rankedDrivers = append(rankedDrivers, driver)
	}

	if len(rankedDrivers) == 0 {
		return nil, fmt.Errorf("no drivers meet acceptance threshold")
	}

	// Sort by score descending
	sort.Slice(rankedDrivers, func(i, j int) bool {
		return rankedDrivers[i].Score > rankedDrivers[j].Score
	})

	return rankedDrivers, nil
}

// SelectBestDriver chooses top driver with confidence calculation
func (mas *MatchingAlgorithmService) SelectBestDriver(
	rankedDrivers []*entities.DriverMatch,
) (*entities.DriverMatch, float64, error) {
	if len(rankedDrivers) == 0 {
		return nil, 0, fmt.Errorf("no ranked drivers available")
	}

	bestDriver := rankedDrivers[0]

	// Calculate confidence based on score distribution
	confidence := calculateConfidence(bestDriver.Score, rankedDrivers)

	return bestDriver, confidence, nil
}

// calculateConfidence determines match confidence (0-1) based on score distribution
func calculateConfidence(topScore float64, rankedDrivers []*entities.DriverMatch) float64 {
	if len(rankedDrivers) < 2 {
		return 1.0 // High confidence if only one option
	}

	secondScore := rankedDrivers[1].Score
	scoreDifference := topScore - secondScore

	// Confidence increases with score gap
	// Max gap is 100 (0 to 100 score range)
	confidence := math.Min(1.0, scoreDifference/50.0)
	return math.Max(0.5, confidence) // Min confidence 0.5
}

// ETACalculatorService handles ETA calculations
type ETACalculatorService struct {
	defaultSpeedKmh float64
}

func NewETACalculatorService(defaultSpeedKmh float64) *ETACalculatorService {
	return &ETACalculatorService{
		defaultSpeedKmh: defaultSpeedKmh,
	}
}

// EstimateETA provides quick ETA estimate based on distance
func (ecs *ETACalculatorService) EstimateETA(
	distanceKm float64,
	trafficFactor float64, // 1.0 = no traffic, 1.5 = 50% traffic delay
) int32 {
	if distanceKm <= 0 {
		return 0
	}

	// Base ETA calculation
	baseDurationHours := distanceKm / ecs.defaultSpeedKmh
	adjustedDurationHours := baseDurationHours * trafficFactor

	// Convert to minutes, add buffer for pickup/dropoff
	etaMinutes := adjustedDurationHours * 60.0
	etaMinutes += 2.0 // 2-minute buffer for pickup

	return int32(math.Ceil(etaMinutes))
}

// SupplyBalancingService handles surge pricing and supply balancing
type SupplyBalancingService struct {
	minSurge float64
	maxSurge float64
}

func NewSupplyBalancingService(minSurge, maxSurge float64) *SupplyBalancingService {
	return &SupplyBalancingService{
		minSurge: minSurge,
		maxSurge: maxSurge,
	}
}

// CalculateSurgeMultiplier determines surge pricing based on supply
func (sbs *SupplyBalancingService) CalculateSurgeMultiplier(
	availableDrivers int,
	pendingRides int,
) float64 {
	if availableDrivers == 0 {
		return sbs.maxSurge
	}

	// Supply-demand ratio
	ratio := float64(pendingRides) / float64(availableDrivers)

	// Exponential surge calculation
	// 1:1 ratio = 1.0 (base), 2:1 = 1.5, 3:1 = 2.0, etc.
	surge := 1.0 + (ratio-1.0)*0.5

	// Clamp to min/max
	if surge < sbs.minSurge {
		surge = sbs.minSurge
	}
	if surge > sbs.maxSurge {
		surge = sbs.maxSurge
	}

	return surge
}

// DriverRankingService provides advanced driver ranking
type DriverRankingService struct {
	matchingService *MatchingAlgorithmService
}

func NewDriverRankingService(ms *MatchingAlgorithmService) *DriverRankingService {
	return &DriverRankingService{
		matchingService: ms,
	}
}

// GetTopDrivers returns top N drivers ranked by score
func (drs *DriverRankingService) GetTopDrivers(
	ctx context.Context,
	drivers []*entities.DriverMatch,
	topN int,
	pickupLat, pickupLng float64,
	acceptanceThreshold float64,
) ([]*entities.DriverMatch, error) {
	rankedDrivers, err := drs.matchingService.RankDrivers(
		ctx, drivers, pickupLat, pickupLng, acceptanceThreshold,
	)
	if err != nil {
		return nil, err
	}

	if topN > len(rankedDrivers) {
		topN = len(rankedDrivers)
	}

	return rankedDrivers[:topN], nil
}
