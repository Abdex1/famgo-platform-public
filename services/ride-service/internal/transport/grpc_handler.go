// services/ride-service/internal/transport/grpc_handler.go
// gRPC Service Handlers

package transport

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"go.uber.org/zap"

	pb "github.com/Abdex1/FamGo-platform/services/ride-service/api/proto/ride"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// GRPCServer implements ride.v1.RideService
type GRPCServer struct {
	pb.UnimplementedRideServiceServer

	createRideHandler     *application.CreateRideHandler
	assignDriverHandler   *application.AssignDriverHandler
	startRideHandler      *application.StartRideHandler
	completeRideHandler   *application.CompleteRideHandler
	cancelRideHandler     *application.CancelRideHandler
	getRideHandler        *application.GetRideHandler
	getPassengerHandler   *application.GetPassengerRidesHandler
	getDriverHandler      *application.GetDriverRidesHandler
	logger                *zap.Logger
}

func NewGRPCServer(
	createRideHandler *application.CreateRideHandler,
	assignDriverHandler *application.AssignDriverHandler,
	startRideHandler *application.StartRideHandler,
	completeRideHandler *application.CompleteRideHandler,
	cancelRideHandler *application.CancelRideHandler,
	getRideHandler *application.GetRideHandler,
	getPassengerHandler *application.GetPassengerRidesHandler,
	getDriverHandler *application.GetDriverRidesHandler,
	logger *zap.Logger,
) *GRPCServer {
	return &GRPCServer{
		createRideHandler:   createRideHandler,
		assignDriverHandler: assignDriverHandler,
		startRideHandler:    startRideHandler,
		completeRideHandler: completeRideHandler,
		cancelRideHandler:   cancelRideHandler,
		getRideHandler:      getRideHandler,
		getPassengerHandler: getPassengerHandler,
		getDriverHandler:    getDriverHandler,
		logger:              logger,
	}
}

// CreateRide creates a new ride
func (s *GRPCServer) CreateRide(ctx context.Context, req *pb.CreateRideRequest) (*pb.CreateRideResponse, error) {
	cmd := application.CreateRideCommand{
		PassengerID: req.PassengerId,
		PickupLat:   req.PickupLat,
		PickupLon:   req.PickupLon,
		DropoffLat:  req.DropoffLat,
		DropoffLon:  req.DropoffLon,
	}

	rideID, err := s.createRideHandler.Handle(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to create ride", zap.Error(err))
		return nil, err
	}

	return &pb.CreateRideResponse{
		RideId: rideID,
		Status: string(domain.RideStatusRequested),
	}, nil
}

// GetRide retrieves ride details
func (s *GRPCServer) GetRide(ctx context.Context, req *pb.GetRideRequest) (*pb.RideResponse, error) {
	ride, err := s.getRideHandler.Handle(ctx, req.RideId)
	if err != nil {
		s.logger.Error("failed to get ride", zap.Error(err))
		return nil, err
	}

	return rideToProto(ride), nil
}

// GetPassengerRides retrieves rides for a passenger
func (s *GRPCServer) GetPassengerRides(ctx context.Context, req *pb.GetPassengerRidesRequest) (*pb.GetPassengerRidesResponse, error) {
	rides, err := s.getPassengerHandler.Handle(ctx, req.PassengerId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("failed to get passenger rides", zap.Error(err))
		return nil, err
	}

	var pbRides []*pb.RideResponse
	for _, ride := range rides {
		pbRides = append(pbRides, rideToProto(&ride))
	}

	return &pb.GetPassengerRidesResponse{
		Rides: pbRides,
		Total: int32(len(pbRides)),
	}, nil
}

// GetDriverRides retrieves rides for a driver
func (s *GRPCServer) GetDriverRides(ctx context.Context, req *pb.GetDriverRidesRequest) (*pb.GetDriverRidesResponse, error) {
	rides, err := s.getDriverHandler.Handle(ctx, req.DriverId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("failed to get driver rides", zap.Error(err))
		return nil, err
	}

	var pbRides []*pb.RideResponse
	for _, ride := range rides {
		pbRides = append(pbRides, rideToProto(&ride))
	}

	return &pb.GetDriverRidesResponse{
		Rides: pbRides,
		Total: int32(len(pbRides)),
	}, nil
}

// AssignDriver assigns a driver to a ride
func (s *GRPCServer) AssignDriver(ctx context.Context, req *pb.AssignDriverRequest) (*pb.AssignDriverResponse, error) {
	cmd := application.AssignDriverCommand{
		RideID:   req.RideId,
		DriverID: req.DriverId,
	}

	err := s.assignDriverHandler.Handle(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to assign driver", zap.Error(err))
		return nil, err
	}

	return &pb.AssignDriverResponse{Success: true}, nil
}

// StartRide marks ride as started
func (s *GRPCServer) StartRide(ctx context.Context, req *pb.StartRideRequest) (*pb.StartRideResponse, error) {
	cmd := application.StartRideCommand{RideID: req.RideId}

	err := s.startRideHandler.Handle(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to start ride", zap.Error(err))
		return nil, err
	}

	return &pb.StartRideResponse{Success: true}, nil
}

// CompleteRide marks ride as completed
func (s *GRPCServer) CompleteRide(ctx context.Context, req *pb.CompleteRideRequest) (*pb.CompleteRideResponse, error) {
	cmd := application.CompleteRideCommand{
		RideID:     req.RideId,
		ActualFare: req.ActualFare,
	}

	err := s.completeRideHandler.Handle(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to complete ride", zap.Error(err))
		return nil, err
	}

	return &pb.CompleteRideResponse{Success: true}, nil
}

// CancelRide cancels a ride
func (s *GRPCServer) CancelRide(ctx context.Context, req *pb.CancelRideRequest) (*pb.CancelRideResponse, error) {
	cmd := application.CancelRideCommand{
		RideID: req.RideId,
		Reason: req.Reason,
	}

	err := s.cancelRideHandler.Handle(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to cancel ride", zap.Error(err))
		return nil, err
	}

	return &pb.CancelRideResponse{Success: true}, nil
}

// Helper: Convert domain.Ride to protobuf RideResponse
func rideToProto(ride *domain.Ride) *pb.RideResponse {
	resp := &pb.RideResponse{
		Id:                 ride.ID,
		PassengerId:        ride.PassengerID,
		DriverId:           ride.DriverID,
		PickupLat:          ride.PickupLat,
		PickupLon:          ride.PickupLon,
		DropoffLat:         ride.DropoffLat,
		DropoffLon:         ride.DropoffLon,
		Status:             string(ride.Status),
		EstimatedFare:      ride.EstimatedFare,
		ActualFare:         ride.ActualFare,
		CancellationReason: ride.CancellationReason,
		CreatedAt:          timestamppb.New(ride.CreatedAt),
		UpdatedAt:          timestamppb.New(ride.UpdatedAt),
	}

	if ride.PickupTime != nil {
		resp.PickupTime = timestamppb.New(*ride.PickupTime)
	}
	if ride.DropoffTime != nil {
		resp.DropoffTime = timestamppb.New(*ride.DropoffTime)
	}

	return resp
}
