// services/ride-service/internal/application/saga.go
// RideCreationSaga - Orchestrates ride creation across services

package application

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// SagaStep represents a single step in the saga
type SagaStep string

const (
	StepCreateRide        SagaStep = "CREATE_RIDE"
	StepRequestDrivers    SagaStep = "REQUEST_DRIVERS"
	StepCalculateFare     SagaStep = "CALCULATE_FARE"
	StepAssignDriver      SagaStep = "ASSIGN_DRIVER"
	StepConfirmPayment    SagaStep = "CONFIRM_PAYMENT"
	StepRideStarted       SagaStep = "RIDE_STARTED"
)

// SagaState represents the current state of a saga
type SagaState struct {
	RideID            string
	PassengerID       string
	PickupLat         float64
	PickupLon         float64
	DropoffLat        float64
	DropoffLon        float64
	CurrentStep       SagaStep
	CompletedSteps    []SagaStep
	FailedStep        SagaStep
	DriverID          string
	CalculatedFare    float32
	Status            string // PENDING, IN_PROGRESS, COMPLETED, FAILED
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// RideCreationSaga orchestrates the ride creation workflow
type RideCreationSaga struct {
	rideRepo      domain.RideRepository
	grpcClients   *RideGRPCClients
	eventPublisher *EventPublisher
	logger        *zap.Logger
}

func NewRideCreationSaga(
	rideRepo domain.RideRepository,
	grpcClients *RideGRPCClients,
	eventPublisher *EventPublisher,
	logger *zap.Logger,
) *RideCreationSaga {
	return &RideCreationSaga{
		rideRepo:      rideRepo,
		grpcClients:   grpcClients,
		eventPublisher: eventPublisher,
		logger:        logger,
	}
}

// ExecuteSaga executes the ride creation saga with compensation
func (s *RideCreationSaga) ExecuteSaga(ctx context.Context, cmd CreateRideCommand) (*SagaState, error) {
	sagaState := &SagaState{
		RideID:      "", // Will be set in step 1
		PassengerID: cmd.PassengerID,
		PickupLat:   cmd.PickupLat,
		PickupLon:   cmd.PickupLon,
		DropoffLat:  cmd.DropoffLat,
		DropoffLon:  cmd.DropoffLon,
		Status:      "IN_PROGRESS",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Step 1: Create Ride
	s.logger.Info("saga step 1: creating ride")
	rideID, err := s.stepCreateRide(ctx, cmd)
	if err != nil {
		s.logger.Error("saga step 1 failed", zap.Error(err))
		sagaState.FailedStep = StepCreateRide
		sagaState.Status = "FAILED"
		return sagaState, err
	}
	sagaState.RideID = rideID
	sagaState.CompletedSteps = append(sagaState.CompletedSteps, StepCreateRide)
	sagaState.CurrentStep = StepRequestDrivers

	// Step 2: Request Drivers
	s.logger.Info("saga step 2: requesting drivers", zap.String("ride_id", rideID))
	drivers, err := s.stepRequestDrivers(ctx, sagaState)
	if err != nil {
		s.logger.Error("saga step 2 failed", zap.Error(err))
		sagaState.FailedStep = StepRequestDrivers
		sagaState.Status = "FAILED"
		s.compensate(ctx, sagaState)
		return sagaState, err
	}
	if len(drivers) == 0 {
		s.logger.Error("no drivers available", zap.String("ride_id", rideID))
		sagaState.FailedStep = StepRequestDrivers
		sagaState.Status = "FAILED"
		s.compensate(ctx, sagaState)
		return sagaState, fmt.Errorf("no drivers available")
	}
	sagaState.CompletedSteps = append(sagaState.CompletedSteps, StepRequestDrivers)
	sagaState.CurrentStep = StepCalculateFare

	// Step 3: Calculate Fare
	s.logger.Info("saga step 3: calculating fare", zap.String("ride_id", rideID))
	fare, err := s.stepCalculateFare(ctx, sagaState)
	if err != nil {
		s.logger.Error("saga step 3 failed", zap.Error(err))
		sagaState.FailedStep = StepCalculateFare
		sagaState.Status = "FAILED"
		s.compensate(ctx, sagaState)
		return sagaState, err
	}
	sagaState.CalculatedFare = fare
	sagaState.CompletedSteps = append(sagaState.CompletedSteps, StepCalculateFare)
	sagaState.CurrentStep = StepAssignDriver

	// Step 4: Assign Driver (use first available)
	s.logger.Info("saga step 4: assigning driver", zap.String("ride_id", rideID))
	driverID := drivers[0]
	err = s.stepAssignDriver(ctx, sagaState, driverID)
	if err != nil {
		s.logger.Error("saga step 4 failed", zap.Error(err))
		sagaState.FailedStep = StepAssignDriver
		sagaState.Status = "FAILED"
		s.compensate(ctx, sagaState)
		return sagaState, err
	}
	sagaState.DriverID = driverID
	sagaState.CompletedSteps = append(sagaState.CompletedSteps, StepAssignDriver)
	sagaState.CurrentStep = StepConfirmPayment

	// Step 5: Confirm Payment
	s.logger.Info("saga step 5: confirming payment", zap.String("ride_id", rideID))
	err = s.stepConfirmPayment(ctx, sagaState)
	if err != nil {
		s.logger.Error("saga step 5 failed", zap.Error(err))
		sagaState.FailedStep = StepConfirmPayment
		sagaState.Status = "FAILED"
		s.compensate(ctx, sagaState)
		return sagaState, err
	}
	sagaState.CompletedSteps = append(sagaState.CompletedSteps, StepConfirmPayment)

	// All steps completed
	sagaState.Status = "COMPLETED"
	sagaState.UpdatedAt = time.Now()

	s.logger.Info("saga completed successfully",
		zap.String("ride_id", rideID),
		zap.String("driver_id", driverID),
		zap.Float32("fare", fare))

	return sagaState, nil
}

// Step 1: Create Ride
func (s *RideCreationSaga) stepCreateRide(ctx context.Context, cmd CreateRideCommand) (string, error) {
	rideID := fmt.Sprintf("ride-%d", time.Now().UnixNano())

	ride := domain.NewRideWithID(
		rideID,
		cmd.PassengerID,
		cmd.PickupLat,
		cmd.PickupLon,
		cmd.DropoffLat,
		cmd.DropoffLon,
	)

	if err := s.rideRepo.Save(ctx, ride); err != nil {
		return "", fmt.Errorf("failed to create ride: %w", err)
	}

	// Publish event
	if err := s.eventPublisher.PublishRideRequested(ctx, rideID, cmd.PassengerID); err != nil {
		s.logger.Warn("failed to publish RideRequested event", zap.Error(err))
		// Don't fail the saga for event publishing issues
	}

	return rideID, nil
}

// Step 2: Request Drivers
func (s *RideCreationSaga) stepRequestDrivers(ctx context.Context, state *SagaState) ([]string, error) {
	// Call dispatch service with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	drivers, err := s.grpcClients.FindDrivers(
		ctx,
		state.RideID,
		state.PickupLat,
		state.PickupLon,
		3, // Find 3 drivers
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find drivers: %w", err)
	}

	return drivers, nil
}

// Step 3: Calculate Fare
func (s *RideCreationSaga) stepCalculateFare(ctx context.Context, state *SagaState) (float32, error) {
	// Call pricing service with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fare, err := s.grpcClients.CalculateFare(
		ctx,
		state.RideID,
		state.PickupLat,
		state.PickupLon,
		state.DropoffLat,
		state.DropoffLon,
		"CAR",
	)

	if err != nil {
		return 0, fmt.Errorf("failed to calculate fare: %w", err)
	}

	return fare, nil
}

// Step 4: Assign Driver
func (s *RideCreationSaga) stepAssignDriver(ctx context.Context, state *SagaState, driverID string) error {
	ride, err := s.rideRepo.GetByID(ctx, state.RideID)
	if err != nil {
		return fmt.Errorf("ride not found: %w", err)
	}

	ride.AssignDriver(driverID)
	ride.SetEstimatedFare(state.CalculatedFare)

	if err := ride.TransitionTo(domain.RideStatusAssigned); err != nil {
		return fmt.Errorf("invalid state transition: %w", err)
	}

	if err := s.rideRepo.Save(ctx, ride); err != nil {
		return fmt.Errorf("failed to save ride: %w", err)
	}

	return nil
}

// Step 5: Confirm Payment
func (s *RideCreationSaga) stepConfirmPayment(ctx context.Context, state *SagaState) error {
	// In production, this would call payment service
	// For now, just verify fare is reasonable
	if state.CalculatedFare <= 0 {
		return fmt.Errorf("invalid fare amount: %f", state.CalculatedFare)
	}

	return nil
}

// Compensate rolls back failed saga steps
func (s *RideCreationSaga) compensate(ctx context.Context, state *SagaState) {
	s.logger.Info("saga compensation started", zap.String("ride_id", state.RideID))

	// Reverse steps in reverse order
	for i := len(state.CompletedSteps) - 1; i >= 0; i-- {
		step := state.CompletedSteps[i]

		switch step {
		case StepCreateRide:
			s.compensateCreateRide(ctx, state)
		case StepAssignDriver:
			s.compensateAssignDriver(ctx, state)
		case StepRequestDrivers, StepCalculateFare, StepConfirmPayment:
			// These are read-only operations, no compensation needed
		}
	}

	s.logger.Info("saga compensation completed", zap.String("ride_id", state.RideID))
}

func (s *RideCreationSaga) compensateCreateRide(ctx context.Context, state *SagaState) {
	s.logger.Info("compensating: cancelling ride", zap.String("ride_id", state.RideID))

	ride, err := s.rideRepo.GetByID(ctx, state.RideID)
	if err != nil {
		s.logger.Error("failed to get ride for compensation", zap.Error(err))
		return
	}

	if err := ride.TransitionTo(domain.RideStatusCancelled); err != nil {
		s.logger.Error("failed to cancel ride", zap.Error(err))
		return
	}

	if err := s.rideRepo.Save(ctx, ride); err != nil {
		s.logger.Error("failed to save cancelled ride", zap.Error(err))
		return
	}

	s.logger.Info("ride cancelled (compensation)")
}

func (s *RideCreationSaga) compensateAssignDriver(ctx context.Context, state *SagaState) {
	s.logger.Info("compensating: releasing driver",
		zap.String("ride_id", state.RideID),
		zap.String("driver_id", state.DriverID))

	// In production, would call dispatch service to release driver
}
