package services

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/valueobjects"
)

// MatchingEngine orchestrates discovery, ranking, and match proposal.
type MatchingEngine struct {
	candidateEngine *DriverCandidateEngine
	matchingService *MatchingService
	topMatchesLimit int
}

type MatchingEngineResult struct {
	PrimaryDriverID string
	ProposedDrivers []string
	Scores          []*valueobjects.MatchScore
}

func NewMatchingEngine(
	candidateEngine *DriverCandidateEngine,
	matchingService *MatchingService,
	topMatchesLimit int,
) *MatchingEngine {
	if topMatchesLimit <= 0 {
		topMatchesLimit = 5
	}
	return &MatchingEngine{
		candidateEngine: candidateEngine,
		matchingService: matchingService,
		topMatchesLimit: topMatchesLimit,
	}
}

// Run executes FindDriversWithinRadius → SortByDistance → score/rank → propose drivers.
func (e *MatchingEngine) Run(
	ctx context.Context,
	request *entities.DispatchRequest,
) (*MatchingEngineResult, error) {
	if request == nil {
		return nil, fmt.Errorf("dispatch request is required")
	}

	candidates, err := e.candidateEngine.FindWithinRadius(
		ctx,
		request.PickupLatitude,
		request.PickupLongitude,
		request.SearchRadius,
		100,
	)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no drivers available within %.2f km", request.SearchRadius)
	}

	sorted := e.candidateEngine.SortByDistance(candidates)
	driverInfos := make([]*DriverInfo, 0, len(sorted))
	for _, candidate := range sorted {
		driverInfos = append(driverInfos, &DriverInfo{
			DriverID:       candidate.DriverID,
			Latitude:       candidate.Latitude,
			Longitude:      candidate.Longitude,
			IsOnline:       candidate.IsOnline,
			AcceptanceRate: candidate.AcceptanceRate,
			Rating:         candidate.Rating,
			Distance:       candidate.DistanceKm,
			ETA:            candidate.ETAMinutes,
		})
	}

	scores, err := e.matchingService.MatchDrivers(
		request.PickupLatitude,
		request.PickupLongitude,
		driverInfos,
		e.topMatchesLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to rank drivers: %w", err)
	}
	if len(scores) == 0 {
		return nil, fmt.Errorf("no eligible drivers after scoring")
	}

	proposed := make([]string, 0, len(scores))
	for _, score := range scores {
		proposed = append(proposed, score.DriverID)
	}

	return &MatchingEngineResult{
		PrimaryDriverID: scores[0].DriverID,
		ProposedDrivers: proposed,
		Scores:          scores,
	}, nil
}

// ExpandAndRetry increases search radius when supply is insufficient.
func (e *MatchingEngine) ExpandAndRetry(
	ctx context.Context,
	request *entities.DispatchRequest,
) (*MatchingEngineResult, error) {
	newRadius := e.matchingService.CalculateOptimalSearchRadius(
		0,
		3,
		request.SearchRadius,
		request.MaxSearchRadius,
	)
	if newRadius <= request.SearchRadius {
		return nil, fmt.Errorf("search radius already at maximum")
	}
	if err := request.ExpandSearchRadius(newRadius); err != nil {
		return nil, err
	}
	return e.Run(ctx, request)
}

// PoolingAwareRun wraps Run with optional pooling eligibility checks.
func (e *MatchingEngine) PoolingAwareRun(
	ctx context.Context,
	request *entities.DispatchRequest,
	poolHook ports.PoolingStrategyHook,
) (*MatchingEngineResult, error) {
	result, err := e.Run(ctx, request)
	if err != nil {
		return nil, err
	}
	if poolHook == nil {
		return result, nil
	}

	canPool, poolErr := poolHook.CanPool(ctx, request.RideID, result.PrimaryDriverID)
	if poolErr != nil {
		return result, nil
	}
	if canPool {
		return result, nil
	}
	return result, nil
}
