// services/ride-service/interfaces/grpc/ride_handler.go
package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/ride-service/proto/ride"
)

type RideHandler struct {
	ride.UnimplementedRideServiceServer
	useCases *usecases.RideUseCases
}

func NewRideHandler(useCases *usecases.RideUseCases) *RideHandler {
	return &RideHandler{useCases: useCases}
}

func (h *RideHandler) CreateRide(ctx context.Context, req *ride.CreateRideRequest) (*ride.CreateRideResponse, error) {
	if req == nil || req.RiderId == "" {
		return nil, status.Error(codes.InvalidArgument, "rider ID is required")
	}

	input := &usecases.CreateRideInput{
		RiderID:        req.RiderId,
		PickupLat:      req.PickupLatitude,
		PickupLng:      req.PickupLongitude,
		DropoffLat:     req.DropoffLatitude,
		DropoffLng:     req.DropoffLongitude,
		PickupAddress:  req.PickupAddress,
		DropoffAddress: req.DropoffAddress,
		PaymentMethod:  req.PaymentMethod,
		PassengerCount: int(req.PassengerCount),
	}

	output, err := h.useCases.CreateRideRequest(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ride.CreateRideResponse{
		RideId:                     output.RideID,
		EstimatedDistance:          output.EstimatedDistance,
		EstimatedDurationMinutes:   output.EstimatedDuration.Minutes(),
		EstimatedFare:              output.EstimatedFare,
	}, nil
}

func (h *RideHandler) AcceptRide(ctx context.Context, req *ride.AcceptRideRequest) (*emptypb.Empty, error) {
	if req == nil || req.RideId == "" || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID and driver ID are required")
	}

	if err := h.useCases.AcceptRide(ctx, &usecases.AcceptRideInput{
		RideID:   req.RideId,
		DriverID: req.DriverId,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *RideHandler) GetRide(ctx context.Context, req *ride.GetRideRequest) (*ride.GetRideResponse, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID is required")
	}

	output, err := h.useCases.GetRideDetails(ctx, req.RideId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	resp := &ride.GetRideResponse{
		RideId:             output.RideID,
		RiderId:            output.RiderID,
		Status:             output.Status,
		EstimatedDistance:  output.EstimatedDistance,
		EstimatedFare:      output.EstimatedFare,
		PickupAddress:      output.PickupAddress,
		DropoffAddress:     output.DropoffAddress,
	}

	if output.DriverID != nil {
		resp.DriverId = *output.DriverID
	}

	if output.ActualFare != nil {
		resp.ActualFare = *output.ActualFare
	}

	if output.DriverRating != nil {
		resp.DriverRating = int32(*output.DriverRating)
	}

	return resp, nil
}

func (h *RideHandler) UpdateRideStatus(ctx context.Context, req *ride.UpdateRideStatusRequest) (*emptypb.Empty, error) {
	if req == nil || req.RideId == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID and status are required")
	}

	if err := h.useCases.UpdateRideStatus(ctx, &usecases.UpdateRideStatusInput{
		RideID: req.RideId,
		Status: req.Status,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *RideHandler) CompleteRide(ctx context.Context, req *ride.CompleteRideRequest) (*emptypb.Empty, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID is required")
	}

	if err := h.useCases.CompleteRide(ctx, &usecases.CompleteRideInput{
		RideID:         req.RideId,
		ActualDistance: req.ActualDistance,
		ActualFare:     req.ActualFare,
		DriverRating:   int(req.DriverRating),
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *RideHandler) CancelRide(ctx context.Context, req *ride.CancelRideRequest) (*emptypb.Empty, error) {
	if req == nil || req.RideId == "" {
		return nil, status.Error(codes.InvalidArgument, "ride ID is required")
	}

	if err := h.useCases.CancelRide(ctx, &usecases.CancelRideInput{
		RideID: req.RideId,
		Reason: req.Reason,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *RideHandler) ListRides(ctx context.Context, req *ride.ListRidesRequest) (*ride.ListRidesResponse, error) {
	// Placeholder implementation
	return &ride.ListRidesResponse{
		Rides: []*ride.RideInfo{},
		Total: 0,
	}, nil
}
