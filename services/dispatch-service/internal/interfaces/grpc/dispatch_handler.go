package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	dispatchv1 "github.com/Abdex1/FamGo-platform/services/dispatch-service/api/proto/v1"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"
)

// InternalDispatchHandler wires gRPC to application use cases (internal adapter layer).
type InternalDispatchHandler struct {
	dispatchv1.UnimplementedDispatchServiceServer
	useCases *usecases.DispatchUseCases
}

func NewInternalDispatchHandler(useCases *usecases.DispatchUseCases) *InternalDispatchHandler {
	return &InternalDispatchHandler{useCases: useCases}
}

func (h *InternalDispatchHandler) MatchRide(ctx context.Context, req *dispatchv1.MatchRideRequest) (*dispatchv1.MatchRideResponse, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID is required")
	}

	input := &usecases.MatchRideInput{
		RideID:            req.RideId,
		RiderID:           req.RiderId,
		PickupLat:         req.PickupLatitude,
		PickupLng:         req.PickupLongitude,
		DropoffLat:        req.DropoffLatitude,
		DropoffLng:        req.DropoffLongitude,
		SearchRadiusKm:    req.SearchRadiusKm,
		MaxSearchRadiusKm: req.MaxSearchRadiusKm,
	}

	result, err := h.useCases.MatchRide(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to match ride: %v", err))
	}

	return &dispatchv1.MatchRideResponse{
		DispatchRequestId: result.DispatchRequestID,
		MatchedDriverId:   result.MatchedDriverID,
		Status:            result.Status,
		ProposedDrivers:   result.ProposedDrivers,
		MatchedAt:         timestamppb.New(result.MatchedAt),
	}, nil
}

func (h *InternalDispatchHandler) GetMatches(ctx context.Context, req *dispatchv1.GetMatchesRequest) (*dispatchv1.GetMatchesResponse, error) {
	if req == nil || req.DispatchRequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "dispatch request ID is required")
	}

	output, err := h.useCases.GetMatches(ctx, &usecases.GetMatchesInput{DispatchRequestID: req.DispatchRequestId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("matches not found: %v", err))
	}

	return &dispatchv1.GetMatchesResponse{
		DispatchRequestId: output.DispatchRequestID,
		Status:            output.Status,
		MatchedDriverId:   stringPtrToString(output.MatchedDriverID),
		ProposedDrivers:   output.ProposedDrivers,
		ExpiresAt:         timestamppb.New(output.ExpiresAt),
	}, nil
}

func (h *InternalDispatchHandler) AcceptMatch(ctx context.Context, req *dispatchv1.AcceptMatchRequest) (*emptypb.Empty, error) {
	if req == nil || req.DispatchRequestId == "" || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "dispatch request ID and driver ID are required")
	}

	if err := h.useCases.AcceptMatch(ctx, &usecases.AcceptMatchInput{
		DispatchRequestID: req.DispatchRequestId,
		DriverID:          req.DriverId,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to accept match: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *InternalDispatchHandler) RejectMatch(ctx context.Context, req *dispatchv1.RejectMatchRequest) (*emptypb.Empty, error) {
	if req == nil || req.DispatchRequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "dispatch request ID is required")
	}

	if err := h.useCases.RejectMatch(ctx, &usecases.RejectMatchInput{
		DispatchRequestID: req.DispatchRequestId,
		Reason:            req.Reason,
		CanRetry:          req.CanRetry,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to reject match: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *InternalDispatchHandler) CancelDispatch(ctx context.Context, req *dispatchv1.CancelDispatchRequest) (*emptypb.Empty, error) {
	if req == nil || req.DispatchRequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "dispatch request ID is required")
	}

	if err := h.useCases.CancelDispatch(ctx, &usecases.CancelDispatchInput{
		DispatchRequestID: req.DispatchRequestId,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to cancel dispatch: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *InternalDispatchHandler) GetDispatchStats(ctx context.Context, req *dispatchv1.GetStatsRequest) (*dispatchv1.GetStatsResponse, error) {
	output, err := h.useCases.GetStats(ctx, &usecases.GetStatsInput{})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get stats: %v", err))
	}

	return &dispatchv1.GetStatsResponse{
		TotalMatches:       int32(output.TotalMatches),
		SuccessfulMatches:  int32(output.SuccessfulMatches),
		FailedMatches:      int32(output.FailedMatches),
		SuccessRate:        output.SuccessRate,
		AverageTimeToMatch: output.AverageTimeToMatch,
	}, nil
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
