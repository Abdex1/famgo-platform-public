package handlers

import (
	"context"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/commands"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/queries"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"
)

// DispatchCommandHandler routes write operations to use cases.
type DispatchCommandHandler struct {
	useCases *usecases.DispatchUseCases
}

func NewDispatchCommandHandler(useCases *usecases.DispatchUseCases) *DispatchCommandHandler {
	return &DispatchCommandHandler{useCases: useCases}
}

func (h *DispatchCommandHandler) MatchRide(ctx context.Context, cmd commands.MatchRideCommand) (*commands.MatchRideResult, error) {
	result, err := h.useCases.MatchRide(ctx, &usecases.MatchRideInput{
		RideID:            cmd.RideID,
		RiderID:           cmd.RiderID,
		PickupLat:         cmd.PickupLat,
		PickupLng:         cmd.PickupLng,
		DropoffLat:        cmd.DropoffLat,
		DropoffLng:        cmd.DropoffLng,
		SearchRadiusKm:    cmd.SearchRadiusKm,
		MaxSearchRadiusKm: cmd.MaxSearchRadiusKm,
		TraceID:           cmd.TraceID,
		CorrelationID:     cmd.CorrelationID,
		RequestID:         cmd.RequestID,
	})
	if err != nil {
		return nil, err
	}
	return &commands.MatchRideResult{
		DispatchRequestID: result.DispatchRequestID,
		MatchedDriverID:   result.MatchedDriverID,
		ProposedDrivers:   result.ProposedDrivers,
		Status:            result.Status,
		MatchedAt:         result.MatchedAt,
	}, nil
}

func (h *DispatchCommandHandler) AcceptMatch(ctx context.Context, cmd commands.AcceptMatchCommand) error {
	return h.useCases.AcceptMatch(ctx, &usecases.AcceptMatchInput{
		DispatchRequestID: cmd.DispatchRequestID,
		DriverID:          cmd.DriverID,
		TraceID:           cmd.TraceID,
		CorrelationID:     cmd.CorrelationID,
		RequestID:         cmd.RequestID,
	})
}

func (h *DispatchCommandHandler) RejectMatch(ctx context.Context, cmd commands.RejectMatchCommand) error {
	return h.useCases.RejectMatch(ctx, &usecases.RejectMatchInput{
		DispatchRequestID: cmd.DispatchRequestID,
		DriverID:          cmd.DriverID,
		Reason:            cmd.Reason,
		CanRetry:          cmd.CanRetry,
	})
}

func (h *DispatchCommandHandler) CancelDispatch(ctx context.Context, cmd commands.CancelDispatchCommand) error {
	return h.useCases.CancelDispatch(ctx, &usecases.CancelDispatchInput{
		DispatchRequestID: cmd.DispatchRequestID,
	})
}

// DispatchQueryHandler routes read operations to use cases.
type DispatchQueryHandler struct {
	useCases *usecases.DispatchUseCases
}

func NewDispatchQueryHandler(useCases *usecases.DispatchUseCases) *DispatchQueryHandler {
	return &DispatchQueryHandler{useCases: useCases}
}

func (h *DispatchQueryHandler) GetMatches(ctx context.Context, query queries.GetMatchesQuery) (*queries.GetMatchesResult, error) {
	output, err := h.useCases.GetMatches(ctx, &usecases.GetMatchesInput{
		DispatchRequestID: query.DispatchRequestID,
	})
	if err != nil {
		return nil, err
	}
	return &queries.GetMatchesResult{
		DispatchRequestID: output.DispatchRequestID,
		Status:            output.Status,
		MatchedDriverID:   output.MatchedDriverID,
		ProposedDrivers:   output.ProposedDrivers,
		ExpiresAt:         output.ExpiresAt,
	}, nil
}

func (h *DispatchQueryHandler) GetStats(ctx context.Context, query queries.GetStatsQuery) (*queries.GetStatsResult, error) {
	output, err := h.useCases.GetStats(ctx, &usecases.GetStatsInput{
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return &queries.GetStatsResult{
		TotalMatches:       output.TotalMatches,
		SuccessfulMatches:  output.SuccessfulMatches,
		FailedMatches:      output.FailedMatches,
		SuccessRate:        output.SuccessRate,
		AverageTimeToMatch: output.AverageTimeToMatch,
	}, nil
}
