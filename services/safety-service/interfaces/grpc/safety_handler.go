// services/safety-service/interfaces/grpc/safety_handler.go
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/safety-service/proto/safety"
)

type SafetyHandler struct {
	safety.UnimplementedSafetyServiceServer
	useCases *usecases.SafetyUseCases
}

func NewSafetyHandler(useCases *usecases.SafetyUseCases) *SafetyHandler {
	return &SafetyHandler{useCases: useCases}
}

func (h *SafetyHandler) InitiateSOS(ctx context.Context, req *safety.InitiateSOSRequest) (*safety.InitiateSOSResponse, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID required")
	}

	output, err := h.useCases.InitiateSOS(ctx, &usecases.InitiateSOSInput{
		RideID:       req.RideId,
		UserID:       req.UserId,
		UserType:     req.UserType,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IncidentType: req.IncidentType,
		Description:  req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to initiate SOS: %v", err))
	}

	return &safety.InitiateSOSResponse{
		IncidentId:      output.IncidentID,
		Status:          output.Status,
		EscalationLevel: output.EscalationLevel,
		InitiatedAt:     timestamppb.Now(),
	}, nil
}

func (h *SafetyHandler) GetIncident(ctx context.Context, req *safety.GetIncidentRequest) (*safety.GetIncidentResponse, error) {
	if req == nil || req.IncidentId == "" {
		return nil, status.Error(codes.InvalidArgument, "incident ID required")
	}

	output, err := h.useCases.GetIncident(ctx, &usecases.GetIncidentInput{IncidentID: req.IncidentId})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("incident not found: %v", err))
	}

	return &safety.GetIncidentResponse{
		IncidentId:      output.IncidentID,
		Status:          output.Status,
		EscalationLevel: output.EscalationLevel,
		InitiatedAt:     timestamppb.Now(),
	}, nil
}

func (h *SafetyHandler) EscalateIncident(ctx context.Context, req *safety.EscalateIncidentRequest) (*emptypb.Empty, error) {
	if req == nil || req.IncidentId == "" {
		return nil, status.Error(codes.InvalidArgument, "incident ID required")
	}

	err := h.useCases.EscalateIncident(ctx, &usecases.EscalateIncidentInput{
		IncidentID:      req.IncidentId,
		EscalationLevel: req.EscalationLevel,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to escalate incident: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *SafetyHandler) ResolveIncident(ctx context.Context, req *safety.ResolveIncidentRequest) (*emptypb.Empty, error) {
	if req == nil || req.IncidentId == "" {
		return nil, status.Error(codes.InvalidArgument, "incident ID required")
	}

	err := h.useCases.ResolveIncident(ctx, &usecases.ResolveIncidentInput{
		IncidentID:     req.IncidentId,
		ResolvedBy:     req.ResolvedBy,
		ResolutionNote: req.ResolutionNotes,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to resolve incident: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (h *SafetyHandler) CancelIncident(ctx context.Context, req *safety.CancelIncidentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
