package grpc

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/gps-service/proto/pb"
)

type GPSHandler struct {
	pb.UnimplementedGPSServiceServer
	updateLocationUC       *usecases.UpdateLocationUseCase
	findNearbyDriversUC    *usecases.FindNearbyDriversUseCase
	driverStatusUC         *usecases.DriverStatusUseCase
}

func NewGPSHandler(
	updateLocationUC *usecases.UpdateLocationUseCase,
	findNearbyDriversUC *usecases.FindNearbyDriversUseCase,
	driverStatusUC *usecases.DriverStatusUseCase,
) *GPSHandler {
	return &GPSHandler{
		updateLocationUC:       updateLocationUC,
		findNearbyDriversUC:    findNearbyDriversUC,
		driverStatusUC:         driverStatusUC,
	}
}

func (h *GPSHandler) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	if req == nil {
		return &pb.UpdateLocationResponse{
			Success:      false,
			ErrorMessage: "request is nil",
		}, nil
	}

	// Call use case
	input := usecases.UpdateLocationInput{
		DriverID:  req.DriverId,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Bearing:   req.Bearing,
		Speed:     req.Speed,
		Provider:  req.Provider,
		Accuracy:  req.Accuracy,
	}

	output, err := h.updateLocationUC.Execute(ctx, input)
	if err != nil {
		return &pb.UpdateLocationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.UpdateLocationResponse{
		LocationId: output.LocationID,
		DriverId:   output.DriverID,
		Timestamp:  getCurrentTimestamp(),
		Success:    true,
	}, nil
}

func (h *GPSHandler) FindNearbyDrivers(ctx context.Context, req *pb.FindNearbyDriversRequest) (*pb.FindNearbyDriversResponse, error) {
	if req == nil {
		return &pb.FindNearbyDriversResponse{
			Success:      false,
			ErrorMessage: "request is nil",
		}, nil
	}

	input := usecases.FindNearbyDriversInput{
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		RadiusKm:   req.RadiusKm,
		OnlineOnly: req.OnlineOnly,
		MaxResults: int(req.MaxResults),
	}

	output, err := h.findNearbyDriversUC.Execute(ctx, input)
	if err != nil {
		return &pb.FindNearbyDriversResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Convert domain drivers to proto
	pbDrivers := make([]*pb.DriverLocation, 0, len(output.Drivers))
	for _, driver := range output.Drivers {
		pbDrivers = append(pbDrivers, &pb.DriverLocation{
			DriverId: driver.DriverID,
			Coordinates: &pb.Coordinates{
				Latitude:  driver.Latitude,
				Longitude: driver.Longitude,
			},
			Bearing:           driver.Bearing,
			Speed:             driver.Speed,
			DistanceFromSearch: driver.Distance,
			LastUpdated:       getCurrentTimestamp(),
			Rating:            driver.Rating,
		})
	}

	return &pb.FindNearbyDriversResponse{
		Drivers:    pbDrivers,
		TotalCount: int32(output.Count),
		Success:    true,
	}, nil
}

func (h *GPSHandler) GetDriverLocation(ctx context.Context, req *pb.GetDriverLocationRequest) (*pb.GetDriverLocationResponse, error) {
	if req == nil || req.DriverId == "" {
		return &pb.GetDriverLocationResponse{
			Success:      false,
			ErrorMessage: "driver_id is required",
		}, nil
	}

	// For now, return placeholder
	// In production, would query repository
	return &pb.GetDriverLocationResponse{
		Location: &pb.DriverLocation{
			DriverId: req.DriverId,
			Coordinates: &pb.Coordinates{
				Latitude:  0,
				Longitude: 0,
			},
		},
		Success: true,
	}, nil
}

func (h *GPSHandler) SetDriverStatus(ctx context.Context, req *pb.SetDriverStatusRequest) (*pb.SetDriverStatusResponse, error) {
	if req == nil {
		return &pb.SetDriverStatusResponse{
			Success:      false,
			ErrorMessage: "request is nil",
		}, nil
	}

	input := usecases.DriverStatusInput{
		DriverID: req.DriverId,
		IsOnline: req.IsOnline,
	}

	output, err := h.driverStatusUC.Execute(ctx, input)
	if err != nil {
		return &pb.SetDriverStatusResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.SetDriverStatusResponse{
		DriverId:  output.DriverID,
		IsOnline:  output.IsOnline,
		Timestamp: getCurrentTimestamp(),
		Success:   output.Updated,
	}, nil
}

func (h *GPSHandler) GetOnlineDrivers(ctx context.Context, req *pb.GetOnlineDriversRequest) (*pb.GetOnlineDriversResponse, error) {
	// Placeholder implementation
	return &pb.GetOnlineDriversResponse{
		Drivers: []*pb.OnlineDriver{},
		Success: true,
	}, nil
}

func (h *GPSHandler) GetLocationHistory(ctx context.Context, req *pb.GetLocationHistoryRequest) (*pb.GetLocationHistoryResponse, error) {
	// Placeholder implementation
	return &pb.GetLocationHistoryResponse{
		History: []*pb.LocationHistoryItem{},
		Success: true,
	}, nil
}

func getCurrentTimestamp() int64 {
	return int64(1000) // Placeholder - use time.Now().UnixMilli() in production
}
