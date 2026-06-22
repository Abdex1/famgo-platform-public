package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
)

// AssignmentEngine assigns drivers and handles reassignment after rejection.
type AssignmentEngine struct {
	matchTimeout time.Duration
	poolHook     ports.PoolingStrategyHook
}

type AssignmentResult struct {
	DispatchRequest *entities.DispatchRequest
	AssignedDriver  string
	ProposedDrivers []string
	AssignedAt      time.Time
}

func NewAssignmentEngine(matchTimeout time.Duration, poolHook ports.PoolingStrategyHook) *AssignmentEngine {
	if matchTimeout <= 0 {
		matchTimeout = 60 * time.Second
	}
	if poolHook == nil {
		poolHook = ports.NoOpPoolingHook{}
	}
	return &AssignmentEngine{
		matchTimeout: matchTimeout,
		poolHook:     poolHook,
	}
}

// AssignFirstAvailable assigns the top-ranked driver from matching output.
func (e *AssignmentEngine) AssignFirstAvailable(
	ctx context.Context,
	request *entities.DispatchRequest,
	matchResult *MatchingEngineResult,
) (*AssignmentResult, error) {
	if request == nil || matchResult == nil {
		return nil, fmt.Errorf("request and match result are required")
	}
	if matchResult.PrimaryDriverID == "" {
		return nil, fmt.Errorf("no primary driver to assign")
	}

	if err := request.StartMatching(); err != nil {
		return nil, err
	}

	canPool, err := e.poolHook.CanPool(ctx, request.RideID, matchResult.PrimaryDriverID)
	if err != nil {
		return nil, fmt.Errorf("pooling eligibility check failed: %w", err)
	}
	_ = canPool

	if err := request.Match(matchResult.PrimaryDriverID, matchResult.ProposedDrivers); err != nil {
		return nil, err
	}

	return &AssignmentResult{
		DispatchRequest: request,
		AssignedDriver:  matchResult.PrimaryDriverID,
		ProposedDrivers: matchResult.ProposedDrivers,
		AssignedAt:      time.Now().UTC(),
	}, nil
}

// Reassign selects the next proposed driver after rejection.
func (e *AssignmentEngine) Reassign(
	request *entities.DispatchRequest,
	rejectedDriverID string,
) (string, error) {
	if request == nil {
		return "", fmt.Errorf("dispatch request is required")
	}

	remaining := make([]string, 0, len(request.ProposedDriverIDs))
	for _, driverID := range request.ProposedDriverIDs {
		if driverID != rejectedDriverID {
			remaining = append(remaining, driverID)
		}
	}
	if len(remaining) == 0 {
		return "", fmt.Errorf("no alternate drivers available for reassignment")
	}

	nextDriver := remaining[0]
	request.Status = entities.StatusMatching
	request.MatchedDriverID = nil
	request.RejectionReason = ""
	request.RejectedAt = nil

	if err := request.Match(nextDriver, remaining); err != nil {
		return "", err
	}
	return nextDriver, nil
}

// MatchTimeout returns the driver response window for a matched request.
func (e *AssignmentEngine) MatchTimeout() time.Duration {
	return e.matchTimeout
}
