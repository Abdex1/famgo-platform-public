// services/ride-service/internal/application/queries.go
// Ride Service Queries and Query Handlers

package application

import (
	"context"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// ===== GET RIDE QUERY =====

type GetRideQuery struct {
	RideID string
}

type GetRideHandler struct {
	rideRepo  domain.RideRepository
	rideCache domain.RideCache
	logger    *zap.Logger
}

func NewGetRideHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	logger *zap.Logger,
) *GetRideHandler {
	return &GetRideHandler{
		rideRepo:  rideRepo,
		rideCache: rideCache,
		logger:    logger,
	}
}

func (h *GetRideHandler) Handle(ctx context.Context, rideID string) (*domain.Ride, error) {
	// Try cache first
	ride, err := h.rideCache.GetRide(ctx, rideID)
	if err == nil && ride != nil {
		h.logger.Debug("ride found in cache", zap.String("rideID", rideID))
		return ride, nil
	}

	// Get from DB
	ride, err = h.rideRepo.GetRide(ctx, rideID)
	if err != nil {
		h.logger.Error("failed to get ride", zap.Error(err))
		return nil, err
	}

	if ride != nil {
		// Cache it
		h.rideCache.SetRide(ctx, ride, 3600)
	}

	return ride, nil
}

// ===== GET RIDES BY PASSENGER QUERY =====

type GetPassengerRidesQuery struct {
	PassengerID string
	Limit       int
	Offset      int
}

type GetPassengerRidesHandler struct {
	rideRepo  domain.RideRepository
	rideCache domain.RideCache
	logger    *zap.Logger
}

func NewGetPassengerRidesHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	logger *zap.Logger,
) *GetPassengerRidesHandler {
	return &GetPassengerRidesHandler{
		rideRepo:  rideRepo,
		rideCache: rideCache,
		logger:    logger,
	}
}

func (h *GetPassengerRidesHandler) Handle(ctx context.Context, passengerID string, limit, offset int) ([]domain.Ride, error) {
	// Default pagination
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	rides, err := h.rideRepo.GetRidesByPassenger(ctx, passengerID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get passenger rides", 
			zap.String("passengerID", passengerID),
			zap.Error(err))
		return nil, err
	}

	h.logger.Info("retrieved passenger rides",
		zap.String("passengerID", passengerID),
		zap.Int("count", len(rides)))

	return rides, nil
}

// ===== GET RIDES BY DRIVER QUERY =====

type GetDriverRidesQuery struct {
	DriverID string
	Limit    int
	Offset   int
}

type GetDriverRidesHandler struct {
	rideRepo  domain.RideRepository
	rideCache domain.RideCache
	logger    *zap.Logger
}

func NewGetDriverRidesHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	logger *zap.Logger,
) *GetDriverRidesHandler {
	return &GetDriverRidesHandler{
		rideRepo:  rideRepo,
		rideCache: rideCache,
		logger:    logger,
	}
}

func (h *GetDriverRidesHandler) Handle(ctx context.Context, driverID string, limit, offset int) ([]domain.Ride, error) {
	// Default pagination
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	rides, err := h.rideRepo.GetRidesByDriver(ctx, driverID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get driver rides",
			zap.String("driverID", driverID),
			zap.Error(err))
		return nil, err
	}

	h.logger.Info("retrieved driver rides",
		zap.String("driverID", driverID),
		zap.Int("count", len(rides)))

	return rides, nil
}

// ===== GET ACTIVE RIDES QUERY =====

type GetActiveRidesQuery struct{}

type GetActiveRidesHandler struct {
	rideRepo  domain.RideRepository
	rideCache domain.RideCache
	logger    *zap.Logger
}

func NewGetActiveRidesHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	logger *zap.Logger,
) *GetActiveRidesHandler {
	return &GetActiveRidesHandler{
		rideRepo:  rideRepo,
		rideCache: rideCache,
		logger:    logger,
	}
}

func (h *GetActiveRidesHandler) Handle(ctx context.Context) ([]domain.Ride, error) {
	// Try cache first
	rides, err := h.rideCache.GetActiveRides(ctx)
	if err == nil && len(rides) > 0 {
		h.logger.Debug("active rides found in cache", zap.Int("count", len(rides)))
		return rides, nil
	}

	// Get from DB
	rides, err = h.rideRepo.GetActiveRides(ctx)
	if err != nil {
		h.logger.Error("failed to get active rides", zap.Error(err))
		return nil, err
	}

	// Cache it
	if len(rides) > 0 {
		h.rideCache.SetActiveRides(ctx, rides, 300) // 5 min cache for active rides
	}

	return rides, nil
}
