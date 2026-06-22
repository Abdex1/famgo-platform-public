// services/ride-service/internal/application/commands.go
// Ride Service Commands and Handlers

package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
	"go.uber.org/zap"
)

// ===== CREATE RIDE COMMAND =====

type CreateRideCommand struct {
	PassengerID string
	PickupLat   float64
	PickupLon   float64
	DropoffLat  float64
	DropoffLon  float64
}

type CreateRideHandler struct {
	rideRepo       domain.RideRepository
	rideCache      domain.RideCache
	rideService    *domain.RideService
	eventPublisher *EventPublisher  // COMPLIANT: Event publishing
	logger         *zap.Logger
}

func NewCreateRideHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	rideService *domain.RideService,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *CreateRideHandler {
	return &CreateRideHandler{
		rideRepo:       rideRepo,
		rideCache:      rideCache,
		rideService:    rideService,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

func (h *CreateRideHandler) Handle(ctx context.Context, cmd CreateRideCommand) (string, error) {
	// Validate locations
	if !h.rideService.ValidateLocation(cmd.PickupLat, cmd.PickupLon) {
		h.logger.Warn("invalid pickup location", zap.Float64("lat", cmd.PickupLat), zap.Float64("lon", cmd.PickupLon))
		return "", errInvalidLocation
	}

	if !h.rideService.ValidateLocation(cmd.DropoffLat, cmd.DropoffLon) {
		h.logger.Warn("invalid dropoff location", zap.Float64("lat", cmd.DropoffLat), zap.Float64("lon", cmd.DropoffLon))
		return "", errInvalidLocation
	}

	// Create ride with ID generated in application layer (Rule 4: domain has ZERO external deps)
	rideID := uuid.New().String()
	ride := domain.NewRideWithID(rideID, cmd.PassengerID, cmd.PickupLat, cmd.PickupLon, cmd.DropoffLat, cmd.DropoffLon)

	// Persist
	if err := h.rideRepo.CreateRide(ctx, ride); err != nil {
		h.logger.Error("failed to create ride", zap.Error(err))
		return "", err
	}

	// Cache
	h.rideCache.SetRide(ctx, ride, 3600)

	// COMPLIANT: Publish event using shared/contracts through packages/event-bus
	correlationID := generateCorrelationID()
	if err := h.eventPublisher.PublishRideRequestedIdempotent(ctx, ride, correlationID); err != nil {
		h.logger.Error("failed to publish ride_requested event", zap.Error(err))
		// Note: Don't fail the command if event publishing fails
		// Events are published asynchronously - use outbox pattern for reliability
	}

	h.logger.Info("ride created", zap.String("rideID", ride.ID), zap.String("passengerID", cmd.PassengerID))

	return ride.ID, nil
}

// ===== ASSIGN DRIVER COMMAND =====

type AssignDriverCommand struct {
	RideID   string
	DriverID string
}

type AssignDriverHandler struct {
	rideRepo       domain.RideRepository
	rideCache      domain.RideCache
	rideService    *domain.RideService
	eventPublisher *EventPublisher
	logger         *zap.Logger
}

func NewAssignDriverHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	rideService *domain.RideService,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *AssignDriverHandler {
	return &AssignDriverHandler{
		rideRepo:       rideRepo,
		rideCache:      rideCache,
		rideService:    rideService,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

func (h *AssignDriverHandler) Handle(ctx context.Context, cmd AssignDriverCommand) error {
	// Get ride
	ride, err := h.rideRepo.GetRide(ctx, cmd.RideID)
	if err != nil {
		h.logger.Error("ride not found", zap.String("rideID", cmd.RideID))
		return err
	}

	// Validate state transition
	if !ride.CanTransitionTo(domain.RideStatusAssigned) {
		h.logger.Warn("cannot assign driver", zap.String("status", string(ride.Status)))
		return errInvalidStateTransition
	}

	// Assign driver
	ride.AssignDriver(cmd.DriverID)
	if err := ride.TransitionTo(domain.RideStatusAssigned); err != nil {
		return err
	}

	// Persist
	if err := h.rideRepo.UpdateRide(ctx, ride); err != nil {
		h.logger.Error("failed to assign driver", zap.Error(err))
		return err
	}

	// Invalidate cache
	h.rideCache.DeleteRide(ctx, cmd.RideID)

	// Publish event
	correlationID := generateCorrelationID()
	if err := h.eventPublisher.PublishRideAssigned(ctx, ride, correlationID); err != nil {
		h.logger.Error("failed to publish ride_assigned event", zap.Error(err))
	}

	h.logger.Info("driver assigned", zap.String("rideID", cmd.RideID), zap.String("driverID", cmd.DriverID))

	return nil
}

// ===== START RIDE COMMAND =====

type StartRideCommand struct {
	RideID string
}

type StartRideHandler struct {
	rideRepo       domain.RideRepository
	rideCache      domain.RideCache
	rideService    *domain.RideService
	eventPublisher *EventPublisher
	logger         *zap.Logger
}

func NewStartRideHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	rideService *domain.RideService,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *StartRideHandler {
	return &StartRideHandler{
		rideRepo:       rideRepo,
		rideCache:      rideCache,
		rideService:    rideService,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

func (h *StartRideHandler) Handle(ctx context.Context, cmd StartRideCommand) error {
	ride, err := h.rideRepo.GetRide(ctx, cmd.RideID)
	if err != nil {
		return err
	}

	// Validate transition
	if !ride.CanTransitionTo(domain.RideStatusStarted) {
		return errInvalidStateTransition
	}

	// Start ride
	ride.StartPickup()
	if err := ride.TransitionTo(domain.RideStatusStarted); err != nil {
		return err
	}

	// Persist
	if err := h.rideRepo.UpdateRide(ctx, ride); err != nil {
		h.logger.Error("failed to start ride", zap.Error(err))
		return err
	}

	h.rideCache.DeleteRide(ctx, cmd.RideID)

	// Publish event
	correlationID := generateCorrelationID()
	if err := h.eventPublisher.PublishRideStarted(ctx, ride, correlationID); err != nil {
		h.logger.Error("failed to publish ride_started event", zap.Error(err))
	}

	h.logger.Info("ride started", zap.String("rideID", cmd.RideID))

	return nil
}

// ===== COMPLETE RIDE COMMAND =====

type CompleteRideCommand struct {
	RideID     string
	ActualFare float32
}

type CompleteRideHandler struct {
	rideRepo       domain.RideRepository
	rideCache      domain.RideCache
	rideService    *domain.RideService
	eventPublisher *EventPublisher
	logger         *zap.Logger
}

func NewCompleteRideHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	rideService *domain.RideService,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *CompleteRideHandler {
	return &CompleteRideHandler{
		rideRepo:       rideRepo,
		rideCache:      rideCache,
		rideService:    rideService,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

func (h *CompleteRideHandler) Handle(ctx context.Context, cmd CompleteRideCommand) error {
	ride, err := h.rideRepo.GetRide(ctx, cmd.RideID)
	if err != nil {
		return err
	}

	// Validate
	if !h.rideService.CanCompleteRide(ride) {
		return errInvalidStateTransition
	}

	// Complete ride
	ride.CompleteRide(cmd.ActualFare)
	if err := ride.TransitionTo(domain.RideStatusCompleted); err != nil {
		return err
	}

	// Persist
	if err := h.rideRepo.UpdateRide(ctx, ride); err != nil {
		h.logger.Error("failed to complete ride", zap.Error(err))
		return err
	}

	h.rideCache.DeleteRide(ctx, cmd.RideID)

	// Publish event
	correlationID := generateCorrelationID()
	if err := h.eventPublisher.PublishRideCompleted(ctx, ride, correlationID); err != nil {
		h.logger.Error("failed to publish ride_completed event", zap.Error(err))
	}

	h.logger.Info("ride completed", zap.String("rideID", cmd.RideID), zap.Float32("fare", cmd.ActualFare))

	return nil
}

// ===== CANCEL RIDE COMMAND =====

type CancelRideCommand struct {
	RideID string
	Reason string
}

type CancelRideHandler struct {
	rideRepo       domain.RideRepository
	rideCache      domain.RideCache
	rideService    *domain.RideService
	eventPublisher *EventPublisher
	logger         *zap.Logger
}

func NewCancelRideHandler(
	rideRepo domain.RideRepository,
	rideCache domain.RideCache,
	rideService *domain.RideService,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *CancelRideHandler {
	return &CancelRideHandler{
		rideRepo:       rideRepo,
		rideCache:      rideCache,
		rideService:    rideService,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

func (h *CancelRideHandler) Handle(ctx context.Context, cmd CancelRideCommand) error {
	ride, err := h.rideRepo.GetRide(ctx, cmd.RideID)
	if err != nil {
		return err
	}

	// Can cancel if not already terminal
	if ride.IsTerminalState() {
		return errInvalidStateTransition
	}

	// Cancel ride
	ride.CancelRide(cmd.Reason)
	if err := ride.TransitionTo(domain.RideStatusCancelled); err != nil {
		return err
	}

	// Persist
	if err := h.rideRepo.UpdateRide(ctx, ride); err != nil {
		h.logger.Error("failed to cancel ride", zap.Error(err))
		return err
	}

	h.rideCache.DeleteRide(ctx, cmd.RideID)

	// Publish event
	correlationID := generateCorrelationID()
	if err := h.eventPublisher.PublishRideCancelled(ctx, ride, correlationID); err != nil {
		h.logger.Error("failed to publish ride_cancelled event", zap.Error(err))
	}

	h.logger.Info("ride cancelled", zap.String("rideID", cmd.RideID), zap.String("reason", cmd.Reason))

	return nil
}

// ===== HELPER FUNCTIONS =====

// generateCorrelationID creates a correlation ID for event tracing
func generateCorrelationID() string {
	return uuid.New().String()
}
