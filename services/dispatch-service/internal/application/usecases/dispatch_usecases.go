// services/dispatch-service/internal/application/usecases/dispatch_usecases.go
package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/repositories"
)

type DispatchUseCases struct {
	repo              *repositories.DispatchRepository
	sessionRepo       *repositories.MatchingSessionRepository
	matchRepo         *repositories.MatchRepository
	matchingEngine    *services.MatchingEngine
	assignmentEngine  *services.AssignmentEngine
	timeoutService    *services.TimeoutService
	events            ports.DispatchEventPublisher
	defaultSearchKm   float64
	maxSearchKm       float64
}

func NewDispatchUseCases(
	repo *repositories.DispatchRepository,
	sessionRepo *repositories.MatchingSessionRepository,
	matchRepo *repositories.MatchRepository,
	matchingEngine *services.MatchingEngine,
	assignmentEngine *services.AssignmentEngine,
	timeoutService *services.TimeoutService,
	events ports.DispatchEventPublisher,
	defaultSearchKm, maxSearchKm float64,
) *DispatchUseCases {
	if events == nil {
		events = eventsNoOp{}
	}
	if defaultSearchKm <= 0 {
		defaultSearchKm = 5.0
	}
	if maxSearchKm <= 0 {
		maxSearchKm = 25.0
	}
	return &DispatchUseCases{
		repo:             repo,
		sessionRepo:      sessionRepo,
		matchRepo:        matchRepo,
		matchingEngine:   matchingEngine,
		assignmentEngine: assignmentEngine,
		timeoutService:   timeoutService,
		events:           events,
		defaultSearchKm:  defaultSearchKm,
		maxSearchKm:      maxSearchKm,
	}
}

type eventsNoOp struct{}

func (eventsNoOp) PublishMatchingStarted(context.Context, ports.MatchingStartedEvent) error { return nil }
func (eventsNoOp) PublishDriverMatched(context.Context, ports.DriverMatchedEvent) error     { return nil }
func (eventsNoOp) PublishDriverAssigned(context.Context, ports.DriverAssignedEvent) error   { return nil }
func (eventsNoOp) PublishMatchingFailed(context.Context, ports.MatchingFailedEvent) error   { return nil }
func (eventsNoOp) PublishMatchingExpired(context.Context, ports.MatchingExpiredEvent) error { return nil }

type MatchRideInput struct {
	RideID            string
	RiderID           string
	PickupLat         float64
	PickupLng         float64
	DropoffLat        float64
	DropoffLng        float64
	SearchRadiusKm    float64
	MaxSearchRadiusKm float64
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type MatchResult struct {
	DispatchRequestID string
	MatchedDriverID   string
	ProposedDrivers   []string
	Status            string
	MatchedAt         time.Time
}

func (uc *DispatchUseCases) MatchRide(ctx context.Context, input *MatchRideInput) (*MatchResult, error) {
	if input == nil || input.RideID == "" {
		return nil, fmt.Errorf("invalid match ride input")
	}

	searchRadius := input.SearchRadiusKm
	if searchRadius <= 0 {
		searchRadius = uc.defaultSearchKm
	}
	maxRadius := input.MaxSearchRadiusKm
	if maxRadius <= 0 {
		maxRadius = uc.maxSearchKm
	}

	dispatchReq, err := entities.NewDispatchRequest(
		input.RideID, input.RiderID,
		input.PickupLat, input.PickupLng,
		input.DropoffLat, input.DropoffLng,
		searchRadius, maxRadius,
		3,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dispatch request: %w", err)
	}
	uc.timeoutService.ApplyExpiry(dispatchReq)

	if err := uc.repo.Create(ctx, dispatchReq); err != nil {
		return nil, fmt.Errorf("failed to persist dispatch request: %w", err)
	}

	session, err := entities.NewMatchingSession(
		dispatchReq.ID,
		dispatchReq.RideID,
		"nearest_available_v1",
		dispatchReq.SearchRadius,
		uc.assignmentEngine.MatchTimeout(),
	)
	if err != nil {
		return nil, err
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to persist matching session: %w", err)
	}

	_ = uc.events.PublishMatchingStarted(ctx, ports.MatchingStartedEvent{
		DispatchRequestID: dispatchReq.ID,
		RideID:            dispatchReq.RideID,
		RiderID:           dispatchReq.RiderID,
		SearchRadiusKm:    dispatchReq.SearchRadius,
		TraceID:           input.TraceID,
		CorrelationID:     input.CorrelationID,
		RequestID:         input.RequestID,
	})

	matchOutput, err := uc.matchingEngine.Run(ctx, dispatchReq)
	if err != nil && dispatchReq.CanRetry() {
		matchOutput, err = uc.matchingEngine.ExpandAndRetry(ctx, dispatchReq)
	}
	if err != nil {
		_ = session.Fail(err.Error())
		_ = uc.sessionRepo.Update(ctx, session)
		dispatchReq.Status = entities.StatusFailed
		_ = uc.repo.Update(ctx, dispatchReq)
		_ = uc.events.PublishMatchingFailed(ctx, ports.MatchingFailedEvent{
			DispatchRequestID: dispatchReq.ID,
			RideID:            dispatchReq.RideID,
			Reason:            err.Error(),
			AttemptCount:      dispatchReq.AttemptCount,
			TraceID:           input.TraceID,
			CorrelationID:     input.CorrelationID,
			RequestID:         input.RequestID,
		})
		return nil, fmt.Errorf("matching failed: %w", err)
	}

	session.SetCandidateCount(len(matchOutput.ProposedDrivers))
	assignment, err := uc.assignmentEngine.AssignFirstAvailable(ctx, dispatchReq, matchOutput)
	if err != nil {
		_ = session.Fail(err.Error())
		_ = uc.sessionRepo.Update(ctx, session)
		return nil, fmt.Errorf("assignment failed: %w", err)
	}

	_ = session.Complete(assignment.AssignedDriver)
	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}
	if err := uc.repo.Update(ctx, assignment.DispatchRequest); err != nil {
		return nil, err
	}

	matchScore := 0.0
	var eta int32
	var distance float64
	if len(matchOutput.Scores) > 0 {
		matchScore = matchOutput.Scores[0].TotalScore
		eta = int32(matchOutput.Scores[0].ETA)
		distance = matchOutput.Scores[0].Distance
	}
	_ = uc.matchRepo.Create(ctx, entities.NewMatchResult(
		session.ID,
		dispatchReq.RideID,
		assignment.AssignedDriver,
		matchScore,
		eta,
		distance,
		0.85,
	))

	_ = uc.events.PublishDriverMatched(ctx, ports.DriverMatchedEvent{
		DispatchRequestID: dispatchReq.ID,
		RideID:            dispatchReq.RideID,
		DriverID:          assignment.AssignedDriver,
		ProposedDrivers:   assignment.ProposedDrivers,
		MatchScore:        matchScore,
		TraceID:           input.TraceID,
		CorrelationID:     input.CorrelationID,
		RequestID:         input.RequestID,
	})

	return &MatchResult{
		DispatchRequestID: dispatchReq.ID,
		MatchedDriverID:   assignment.AssignedDriver,
		ProposedDrivers:   assignment.ProposedDrivers,
		Status:            string(assignment.DispatchRequest.Status),
		MatchedAt:         assignment.AssignedAt,
	}, nil
}

type GetMatchesInput struct {
	DispatchRequestID string
}

type GetMatchesOutput struct {
	DispatchRequestID string
	Status            string
	MatchedDriverID   *string
	ProposedDrivers   []string
	ExpiresAt         time.Time
}

func (uc *DispatchUseCases) GetMatches(ctx context.Context, input *GetMatchesInput) (*GetMatchesOutput, error) {
	if input == nil || input.DispatchRequestID == "" {
		return nil, fmt.Errorf("invalid get matches input")
	}

	dispatchReq, err := uc.repo.GetByID(ctx, input.DispatchRequestID)
	if err != nil {
		return nil, fmt.Errorf("dispatch request not found: %w", err)
	}
	if err := uc.timeoutService.CheckAndExpire(dispatchReq); err == nil && dispatchReq.Status == entities.StatusExpired {
		_ = uc.repo.Update(ctx, dispatchReq)
	}

	return &GetMatchesOutput{
		DispatchRequestID: dispatchReq.ID,
		Status:            string(dispatchReq.Status),
		MatchedDriverID:   dispatchReq.MatchedDriverID,
		ProposedDrivers:   dispatchReq.ProposedDriverIDs,
		ExpiresAt:         dispatchReq.ExpiryTime,
	}, nil
}

type AcceptMatchInput struct {
	DispatchRequestID string
	DriverID          string
	TraceID           string
	CorrelationID     string
	RequestID         string
}

func (uc *DispatchUseCases) AcceptMatch(ctx context.Context, input *AcceptMatchInput) error {
	if input == nil || input.DispatchRequestID == "" || input.DriverID == "" {
		return fmt.Errorf("invalid accept match input")
	}

	dispatchReq, err := uc.repo.GetByID(ctx, input.DispatchRequestID)
	if err != nil {
		return fmt.Errorf("dispatch request not found: %w", err)
	}

	if err := dispatchReq.Accept(input.DriverID); err != nil {
		return fmt.Errorf("failed to accept match: %w", err)
	}

	if err := uc.repo.Update(ctx, dispatchReq); err != nil {
		return fmt.Errorf("failed to update dispatch request: %w", err)
	}

	return uc.events.PublishDriverAssigned(ctx, ports.DriverAssignedEvent{
		DispatchRequestID: dispatchReq.ID,
		RideID:            dispatchReq.RideID,
		DriverID:          input.DriverID,
		TraceID:           input.TraceID,
		CorrelationID:     input.CorrelationID,
		RequestID:         input.RequestID,
	})
}

type RejectMatchInput struct {
	DispatchRequestID string
	Reason            string
	CanRetry          bool
	DriverID          string
}

func (uc *DispatchUseCases) RejectMatch(ctx context.Context, input *RejectMatchInput) error {
	if input == nil || input.DispatchRequestID == "" {
		return fmt.Errorf("invalid reject match input")
	}

	dispatchReq, err := uc.repo.GetByID(ctx, input.DispatchRequestID)
	if err != nil {
		return fmt.Errorf("dispatch request not found: %w", err)
	}

	rejectedDriver := input.DriverID
	if rejectedDriver == "" && dispatchReq.MatchedDriverID != nil {
		rejectedDriver = *dispatchReq.MatchedDriverID
	}

	if err := dispatchReq.Reject(input.Reason); err != nil {
		return fmt.Errorf("failed to reject match: %w", err)
	}

	if input.CanRetry && uc.timeoutService.ShouldRetryAfterTimeout(dispatchReq) {
		if _, reassignErr := uc.assignmentEngine.Reassign(dispatchReq, rejectedDriver); reassignErr != nil {
			if err := dispatchReq.RetryMatching(); err != nil {
				return fmt.Errorf("failed to retry matching: %w", err)
			}
			dispatchReq.Status = entities.StatusFailed
		}
	}

	if err := uc.repo.Update(ctx, dispatchReq); err != nil {
		return fmt.Errorf("failed to update dispatch request: %w", err)
	}
	return nil
}

type GetStatsInput struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetStatsOutput struct {
	TotalMatches       int
	SuccessfulMatches  int
	FailedMatches      int
	SuccessRate        float64
	AverageTimeToMatch float64
}

func (uc *DispatchUseCases) GetStats(ctx context.Context, input *GetStatsInput) (*GetStatsOutput, error) {
	stats, err := uc.repo.GetDispatchStats(ctx, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}
	return &GetStatsOutput{
		TotalMatches:       stats.TotalMatches,
		SuccessfulMatches:  stats.SuccessfulMatches,
		FailedMatches:      stats.FailedMatches,
		SuccessRate:        stats.SuccessRate,
		AverageTimeToMatch: stats.AverageTimeToMatch,
	}, nil
}

type CancelDispatchInput struct {
	DispatchRequestID string
}

func (uc *DispatchUseCases) CancelDispatch(ctx context.Context, input *CancelDispatchInput) error {
	if input == nil || input.DispatchRequestID == "" {
		return fmt.Errorf("invalid cancel dispatch input")
	}

	dispatchReq, err := uc.repo.GetByID(ctx, input.DispatchRequestID)
	if err != nil {
		return fmt.Errorf("dispatch request not found: %w", err)
	}

	if err := dispatchReq.Cancel(); err != nil {
		return fmt.Errorf("failed to cancel: %w", err)
	}

	if err := uc.repo.Update(ctx, dispatchReq); err != nil {
		return fmt.Errorf("failed to update dispatch request: %w", err)
	}

	if session, err := uc.sessionRepo.GetByDispatchRequestID(ctx, dispatchReq.ID); err == nil {
		session.Cancel()
		_ = uc.sessionRepo.Update(ctx, session)
	}
	return nil
}
