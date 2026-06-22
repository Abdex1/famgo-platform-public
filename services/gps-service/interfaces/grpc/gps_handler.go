// services/gps-service/interfaces/grpc/gps_handler.go
// gRPC handler for GPS service

package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/gps-service/proto/gps"
)

// GPSHandler implements the GPSService gRPC service
type GPSHandler struct {
	gps.UnimplementedGPSServiceServer
	useCases *usecases.LocationUseCases
}

// NewGPSHandler creates a new GPS handler
func NewGPSHandler(useCases *usecases.LocationUseCases) *GPSHandler {
	return &GPSHandler{
		useCases: useCases,
	}
}

// UpdateLocation updates a driver's location
func (h *GPSHandler) UpdateLocation(
	ctx context.Context,
	req *gps.UpdateLocationRequest,
) (*gps.UpdateLocationResponse, error) {
	if req == nil || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "driver ID is required")
	}

	input := &usecases.UpdateLocationInput{
		DriverID:  req.DriverId,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Altitude:  req.Altitude,
		Accuracy:  req.Accuracy,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Timestamp: req.Timestamp,
	}

	output, err := h.useCases.UpdateDriverLocation(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update location: %v", err))
	}

	return &gps.UpdateLocationResponse{
		Success:            output.Success,
		DriverId:           output.DriverID,
		Distance:           output.Distance,
		EtaMinutes:         output.ETAMinutes,
		QualityIssues:      output.QualityIssues,
		ConsecutiveFails:   int32(output.ConsecutiveFails),
	}, nil
}

// FindNearbyDrivers finds drivers within a radius
func (h *GPSHandler) FindNearbyDrivers(
	ctx context.Context,
	req *gps.FindNearbyDriversRequest,
) (*gps.FindNearbyDriversResponse, error) {
	if req == nil || req.RadiusKm <= 0 {
		return nil, status.Error(codes.InvalidArgument, "valid radius is required")
	}

	input := &usecases.FindNearbyDriversInput{
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusKm:     req.RadiusKm,
		Limit:        int(req.Limit),
		BaseSpeedKmH: req.BaseSpeedKmh,
		OnlyOnline:   req.OnlyOnline,
	}

	if input.Limit <= 0 {
		input.Limit = 50
	}
	if input.BaseSpeedKmH <= 0 {
		input.BaseSpeedKmH = 40
	}

	results, err := h.useCases.FindNearbyDrivers(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find nearby drivers: %v", err))
	}

	drivers := make([]*gps.NearbyDriver, len(results))
	for i, result := range results {
		drivers[i] = &gps.NearbyDriver{
			DriverId:       result.DriverID,
			Distance:       result.Distance,
			EtaMinutes:     result.ETAMinutes,
			Bearing:        result.Bearing,
			Latitude:       result.Latitude,
			Longitude:      result.Longitude,
			IsOnline:       result.IsOnline,
			AcceptanceRate: result.AcceptanceRate,
			Rating:         result.Rating,
			RidesCompleted: int32(result.RidesCompleted),
		}
	}

	return &gps.FindNearbyDriversResponse{
		Drivers: drivers,
		Total:   int32(len(drivers)),
	}, nil
}

// GetDriverLocation retrieves current driver location
func (h *GPSHandler) GetDriverLocation(
	ctx context.Context,
	req *gps.GetDriverLocationRequest,
) (*gps.GetDriverLocationResponse, error) {
	if req == nil || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "driver ID is required")
	}

	output, err := h.useCases.GetDriverLocation(ctx, req.DriverId)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("driver location not found: %v", err))
	}

	return &gps.GetDriverLocationResponse{
		DriverId:       output.DriverID,
		Latitude:       output.Latitude,
		Longitude:      output.Longitude,
		Altitude:       output.Altitude,
		Accuracy:       output.Accuracy,
		Speed:          output.Speed,
		Heading:        output.Heading,
		Status:         output.Status,
		IsOnline:       output.IsOnline,
		LastUpdateAt:   timestamppb.New(output.LastUpdateAt),
		AcceptanceRate: output.AcceptanceRate,
		Rating:         output.Rating,
		RidesCompleted: int32(output.RidesCompleted),
	}, nil
}

// UpdateDriverStatus updates driver online/offline status
func (h *GPSHandler) UpdateDriverStatus(
	ctx context.Context,
	req *gps.UpdateDriverStatusRequest,
) (*emptypb.Empty, error) {
	if req == nil || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "driver ID is required")
	}

	input := &usecases.UpdateDriverStatusInput{
		DriverID: req.DriverId,
		Status:   req.Status,
	}

	if err := h.useCases.UpdateDriverStatus(ctx, input); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update driver status: %v", err))
	}

	return &emptypb.Empty{}, nil
}

// BulkUpdateLocations updates multiple driver locations
func (h *GPSHandler) BulkUpdateLocations(
	ctx context.Context,
	req *gps.BulkUpdateLocationsRequest,
) (*gps.BulkUpdateLocationsResponse, error) {
	if req == nil || len(req.Locations) == 0 {
		return nil, status.Error(codes.InvalidArgument, "locations are required")
	}

	// Convert proto locations to use case inputs
	locations := make([]usecases.UpdateLocationInput, len(req.Locations))
	for i, loc := range req.Locations {
		locations[i] = usecases.UpdateLocationInput{
			DriverID:  loc.DriverId,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Altitude:  loc.Altitude,
			Accuracy:  loc.Accuracy,
			Speed:     loc.Speed,
			Heading:   loc.Heading,
			Timestamp: loc.Timestamp,
		}
	}

	input := &usecases.BulkUpdateLocationsInput{
		Locations: locations,
	}

	output, err := h.useCases.BulkUpdateLocations(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to bulk update locations: %v", err))
	}

	return &gps.BulkUpdateLocationsResponse{
		Processed:     int32(output.Processed),
		Failed:        int32(output.Failed),
		FailedDrivers: output.FailedDrivers,
	}, nil
}

// GetOnlineDriversCount returns count of online drivers
func (h *GPSHandler) GetOnlineDriversCount(
	ctx context.Context,
	_ *emptypb.Empty,
) (*gps.GetOnlineDriversCountResponse, error) {
	onlineDrivers, err := h.useCases.driverTrackingStore.GetOnlineDrivers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get online drivers count: %v", err))
	}

	return &gps.GetOnlineDriversCountResponse{
		Count: int64(len(onlineDrivers)),
	}, nil
}
