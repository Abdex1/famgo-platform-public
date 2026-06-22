// services/fraud-service/interfaces/grpc/fraud_handler.go
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/fraud-service/proto/fraud"
)

type FraudHandler struct {
	fraud.UnimplementedFraudServiceServer
	useCases *usecases.FraudUseCases
}

func NewFraudHandler(useCases *usecases.FraudUseCases) *FraudHandler {
	return &FraudHandler{useCases: useCases}
}

func (h *FraudHandler) CheckRide(ctx context.Context, req *fraud.CheckRideRequest) (*fraud.CheckRideResponse, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID required")
	}

	output, err := h.useCases.CheckRide(ctx, &usecases.CheckRideInput{
		RideID:    req.RideId,
		UserID:    req.UserId,
		UserType:  req.UserType,
		Amount:    req.Amount,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check ride: %v", err))
	}

	return &fraud.CheckRideResponse{
		CheckId:    output.CheckID,
		RiskScore:  output.RiskScore,
		RiskLevel:  output.RiskLevel,
		Action:     output.Action,
		Flags:      output.FlagsTriggered,
		CheckedAt:  timestamppb.Now(),
	}, nil
}

func (h *FraudHandler) GetFraudCheck(ctx context.Context, req *fraud.GetFraudCheckRequest) (*fraud.GetFraudCheckResponse, error) {
	if req == nil || req.CheckId == "" {
		return nil, status.Error(codes.InvalidArgument, "check ID required")
	}

	output, err := h.useCases.GetCheck(ctx, &usecases.GetCheckInput{CheckID: req.CheckId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("fraud check not found: %v", err))
	}

	return &fraud.GetFraudCheckResponse{
		CheckId:    output.CheckID,
		RiskScore:  output.RiskScore,
		RiskLevel:  output.RiskLevel,
		Flags:      output.FlagsTriggered,
		IsReview:   true,
		CreatedAt:  timestamppb.Now(),
	}, nil
}

func (h *FraudHandler) ReviewCheck(ctx context.Context, req *fraud.ReviewCheckRequest) (*emptypb.Empty, error) {
	if req == nil || req.CheckId == "" {
		return nil, status.Error(codes.InvalidArgument, "check ID required")
	}

	err := h.useCases.ReviewCheck(ctx, &usecases.ReviewCheckInput{
		CheckID:    req.CheckId,
		ReviewedBy: req.ReviewedBy,
		Reason:     req.Reason,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to review check: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *FraudHandler) OverrideCheck(ctx context.Context, req *fraud.OverrideCheckRequest) (*emptypb.Empty, error) {
	if req == nil || req.CheckId == "" {
		return nil, status.Error(codes.InvalidArgument, "check ID required")
	}

	err := h.useCases.OverrideCheck(ctx, &usecases.OverrideCheckInput{
		CheckID: req.CheckId,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to override check: %v", err))
	}

	return &emptypb.Empty{}, nil
}
