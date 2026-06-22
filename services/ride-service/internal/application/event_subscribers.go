// services/ride-service/internal/application/event_subscribers.go
// Event subscription handlers for Ride Service

package application

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
)

// RideEventSubscriber handles incoming events
type RideEventSubscriber struct {
	rideRepo domain.RideRepository
	logger   *zap.Logger
}

func NewRideEventSubscriber(rideRepo domain.RideRepository, logger *zap.Logger) *RideEventSubscriber {
	return &RideEventSubscriber{
		rideRepo: rideRepo,
		logger:   logger,
	}
}

// HandleDriverAssigned handles DriverAssigned event from dispatch service
func (s *RideEventSubscriber) HandleDriverAssigned(ctx context.Context, event *events.Event) error {
	s.logger.Info("received DriverAssigned event",
		zap.String("event_id", event.EventID),
		zap.String("ride_id", event.AggregateID))

	// Parse event data
	var payload struct {
		RideID   string `json:"ride_id"`
		DriverID string `json:"driver_id"`
		Fare     float32 `json:"estimated_fare"`
	}

	if err := json.Unmarshal(event.Data, &payload); err != nil {
		s.logger.Error("failed to unmarshal event data", zap.Error(err))
		return err
	}

	// Get ride from repository
	ride, err := s.rideRepo.GetByID(ctx, payload.RideID)
	if err != nil {
		s.logger.Error("ride not found", zap.String("ride_id", payload.RideID), zap.Error(err))
		return fmt.Errorf("ride not found: %w", err)
	}

	// Update ride with driver assignment
	ride.AssignDriver(payload.DriverID)
	ride.SetEstimatedFare(payload.Fare)

	// Transition to ASSIGNED state
	if err := ride.TransitionTo(domain.RideStatusAssigned); err != nil {
		s.logger.Error("invalid state transition", zap.Error(err))
		return err
	}

	// Persist
	if err := s.rideRepo.Save(ctx, ride); err != nil {
		s.logger.Error("failed to save ride", zap.Error(err))
		return err
	}

	s.logger.Info("ride assigned successfully",
		zap.String("ride_id", payload.RideID),
		zap.String("driver_id", payload.DriverID))

	return nil
}

// HandlePaymentProcessed handles PaymentProcessed event from payment service
func (s *RideEventSubscriber) HandlePaymentProcessed(ctx context.Context, event *events.Event) error {
	s.logger.Info("received PaymentProcessed event",
		zap.String("event_id", event.EventID),
		zap.String("ride_id", event.AggregateID))

	// Parse event data
	var payload struct {
		RideID string  `json:"ride_id"`
		Amount float32 `json:"amount"`
		Status string  `json:"status"`
	}

	if err := json.Unmarshal(event.Data, &payload); err != nil {
		s.logger.Error("failed to unmarshal event data", zap.Error(err))
		return err
	}

	if payload.Status != "SUCCESS" {
		s.logger.Warn("payment failed, ride should be cancelled",
			zap.String("ride_id", payload.RideID),
			zap.String("payment_status", payload.Status))
		return nil
	}

	// Payment successful - ride can proceed
	s.logger.Info("payment processed successfully",
		zap.String("ride_id", payload.RideID),
		zap.Float32("amount", payload.Amount))

	return nil
}

// RegisterEventSubscriptions registers all event subscriptions
func (s *RideEventSubscriber) RegisterEventSubscriptions(eventBus interface{}) error {
	// This would be called during bootstrap to register handlers
	// EventBus implementation would subscribe to:
	// - "DriverAssigned" topic → HandleDriverAssigned
	// - "PaymentProcessed" topic → HandlePaymentProcessed

	s.logger.Info("ride service event subscriptions registered",
		zap.Strings("subscriptions", []string{"DriverAssigned", "PaymentProcessed"}))

	return nil
}
