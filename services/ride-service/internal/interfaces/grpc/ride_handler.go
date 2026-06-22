package grpc

import (
	"context"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/ride-service/proto/pb"
)

type RideHandler struct {
	pb.UnimplementedRideServiceServer
	createRideUC   *usecases.CreateRideUseCase
	acceptRideUC   *usecases.AcceptRideUseCase
	startRideUC    *usecases.StartRideUseCase
	completeRideUC *usecases.CompleteRideUseCase
	cancelRideUC   *usecases.CancelRideUseCase
}

func NewRideHandler(
	createRideUC *usecases.CreateRideUseCase,
	acceptRideUC *usecases.AcceptRideUseCase,
	startRideUC *usecases.StartRideUseCase,
	completeRideUC *usecases.CompleteRideUseCase,
	cancelRideUC *usecases.CancelRideUseCase,
) *RideHandler {
	return &RideHandler{
		createRideUC:   createRideUC,
		acceptRideUC:   acceptRideUC,
		startRideUC:    startRideUC,
		completeRideUC: completeRideUC,
		cancelRideUC:   cancelRideUC,
	}
}

func (h *RideHandler) CreateRide(ctx context.Context, req *pb.CreateRideRequest) (*pb.CreateRideResponse, error) {
	if req == nil {
		return &pb.CreateRideResponse{Success: false, ErrorMessage: "request is nil"}, nil
	}

	input := usecases.CreateRideInput{
		RiderID:           req.RiderId,
		PickupLat:         req.Pickup.Latitude,
		PickupLng:         req.Pickup.Longitude,
		DropoffLat:        req.Dropoff.Latitude,
		DropoffLng:        req.Dropoff.Longitude,
		EstimatedDistance: req.EstimatedDistance,
		EstimatedDuration: req.EstimatedDuration,
		PaymentMethod:     req.PaymentMethod,
	}

	output, err := h.createRideUC.Execute(ctx, input)
	if err != nil {
		return &pb.CreateRideResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	return &pb.CreateRideResponse{
		RideId:            output.RideID,
		Status:            output.Status,
		EstimatedDistance: output.EstimatedDistance,
		EstimatedDuration: output.EstimatedDuration,
		Success:           true,
	}, nil
}

func (h *RideHandler) AcceptRide(ctx context.Context, req *pb.AcceptRideRequest) (*pb.AcceptRideResponse, error) {
	if req == nil {
		return &pb.AcceptRideResponse{Success: false, ErrorMessage: "request is nil"}, nil
	}

	input := usecases.AcceptRideInput{
		RideID:   req.RideId,
		DriverID: "",
	}

	output, err := h.acceptRideUC.Execute(ctx, input)
	if err != nil {
		return &pb.AcceptRideResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	return &pb.AcceptRideResponse{
		RideId:  output.RideID,
		Status:  output.Status,
		Success: true,
	}, nil
}

func (h *RideHandler) StartRide(ctx context.Context, req *pb.StartRideRequest) (*pb.StartRideResponse, error) {
	if req == nil {
		return &pb.StartRideResponse{Success: false, ErrorMessage: "request is nil"}, nil
	}

	input := usecases.StartRideInput{
		RideID: req.RideId,
	}

	output, err := h.startRideUC.Execute(ctx, input)
	if err != nil {
		return &pb.StartRideResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	return &pb.StartRideResponse{
		RideId:    output.RideID,
		Status:    output.Status,
		StartedAt: getCurrentTimestamp(),
		Success:   true,
	}, nil
}

func (h *RideHandler) CompleteRide(ctx context.Context, req *pb.CompleteRideRequest) (*pb.CompleteRideResponse, error) {
	if req == nil {
		return &pb.CompleteRideResponse{Success: false, ErrorMessage: "request is nil"}, nil
	}

	input := usecases.CompleteRideInput{
		RideID:          req.RideId,
		ActualDistance:  req.ActualDistance,
		ActualDuration:  req.ActualDuration,
		BaseFare:        req.BaseFare,
		DistanceFare:    req.DistanceFare,
		TimeFare:        req.TimeFare,
		SurgeMultiplier: req.SurgeMultiplier,
	}

	output, err := h.completeRideUC.Execute(ctx, input)
	if err != nil {
		return &pb.CompleteRideResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	return &pb.CompleteRideResponse{
		RideId:      output.RideID,
		Status:      output.Status,
		TotalFare:   output.TotalFare,
		CompletedAt: getCurrentTimestamp(),
		Success:     true,
	}, nil
}

func (h *RideHandler) CancelRide(ctx context.Context, req *pb.CancelRideRequest) (*pb.CancelRideResponse, error) {
	if req == nil {
		return &pb.CancelRideResponse{Success: false, ErrorMessage: "request is nil"}, nil
	}

	input := usecases.CancelRideInput{
		RideID: req.RideId,
		Reason: req.Reason,
		CancellationFeeAmount: req.CancellationFee,
	}

	output, err := h.cancelRideUC.Execute(ctx, input)
	if err != nil {
		return &pb.CancelRideResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	return &pb.CancelRideResponse{
		RideId:            output.RideID,
		Status:            output.Status,
		CancellationFee:   output.CancellationFeeAmount,
		CancelledAt:       getCurrentTimestamp(),
		Success:           true,
	}, nil
}

func (h *RideHandler) AssignDriver(ctx context.Context, req *pb.AssignDriverRequest) (*pb.AssignDriverResponse, error) {
	return &pb.AssignDriverResponse{Success: true}, nil
}

func (h *RideHandler) GetRide(ctx context.Context, req *pb.GetRideRequest) (*pb.GetRideResponse, error) {
	return &pb.GetRideResponse{Success: true}, nil
}

func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
