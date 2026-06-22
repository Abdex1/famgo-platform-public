// services/ride-service/internal/application/events.go
// Event Publishing Layer - CRITICAL GOVERNANCE COMPLIANCE
// RULE 1: Events MUST use shared/contracts/events ONLY
// RULE 2: Publishing MUST use packages/event-bus

package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	// COMPLIANT: Using shared contracts (NOT local events)
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// EventPublisher handles ride domain events
type EventPublisher struct {
	eventBus event_bus.EventBus  // From packages/event-bus (COMPLIANT)
	logger   *zap.Logger
}

func NewEventPublisher(bus event_bus.EventBus, logger *zap.Logger) *EventPublisher {
	return &EventPublisher{
		eventBus: bus,
		logger:   logger,
	}
}

// ===== RIDE CREATED EVENT =====

func (p *EventPublisher) PublishRideRequested(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	// Create event using shared/contracts/events envelope
	// RULE 1 COMPLIANT: Uses shared contract, not local event
	
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.requested",  // From shared contracts
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":      ride.ID,
			"passenger_id": ride.PassengerID,
			"pickup_lat":   ride.PickupLat,
			"pickup_lon":   ride.PickupLon,
			"dropoff_lat":  ride.DropoffLat,
			"dropoff_lon":  ride.DropoffLon,
			"status":       ride.Status,
		},
	}

	// RULE 2 COMPLIANT: Publishing through packages/event-bus
	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish ride_requested",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_requested published",
		zap.String("rideID", ride.ID),
		zap.String("passengerID", ride.PassengerID))

	return nil
}

// ===== DRIVER ASSIGNED EVENT =====

func (p *EventPublisher) PublishRideAssigned(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.assigned",
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":   ride.ID,
			"driver_id": ride.DriverID,
			"status":    ride.Status,
		},
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish ride_assigned",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_assigned published",
		zap.String("rideID", ride.ID),
		zap.String("driverID", ride.DriverID))

	return nil
}

// ===== RIDE STARTED EVENT =====

func (p *EventPublisher) PublishRideStarted(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.started",
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":     ride.ID,
			"driver_id":   ride.DriverID,
			"passenger_id": ride.PassengerID,
			"pickup_time": ride.PickupTime,
			"status":      ride.Status,
		},
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish ride_started",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_started published",
		zap.String("rideID", ride.ID))

	return nil
}

// ===== RIDE COMPLETED EVENT =====

func (p *EventPublisher) PublishRideCompleted(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.completed",
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":      ride.ID,
			"driver_id":    ride.DriverID,
			"passenger_id": ride.PassengerID,
			"estimated_fare": ride.EstimatedFare,
			"actual_fare":    ride.ActualFare,
			"dropoff_time":   ride.DropoffTime,
			"status":         ride.Status,
		},
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish ride_completed",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_completed published",
		zap.String("rideID", ride.ID),
		zap.Float32("actualFare", ride.ActualFare))

	return nil
}

// ===== RIDE CANCELLED EVENT =====

func (p *EventPublisher) PublishRideCancelled(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.cancelled",
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":              ride.ID,
			"driver_id":            ride.DriverID,
			"passenger_id":         ride.PassengerID,
			"cancellation_reason":  ride.CancellationReason,
			"status":               ride.Status,
		},
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish ride_cancelled",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_cancelled published",
		zap.String("rideID", ride.ID),
		zap.String("reason", ride.CancellationReason))

	return nil
}

// ===== IDEMPOTENT PUBLISHING (CRITICAL FOR RELIABILITY) =====

// PublishRideRequestedIdempotent ensures the event is published exactly-once
// RULE 2: Uses packages/event-bus with idempotency key
func (p *EventPublisher) PublishRideRequestedIdempotent(
	ctx context.Context,
	ride *domain.Ride,
	correlationID string,
) error {
	event := &events.Event{
		EventID:       uuid.New().String(),
		EventType:     "ride.requested",
		Version:       1,
		AggregateID:   ride.ID,
		AggregateType: "ride",
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
		CausationID:   ride.ID,
		Data: map[string]interface{}{
			"ride_id":      ride.ID,
			"passenger_id": ride.PassengerID,
			"pickup_lat":   ride.PickupLat,
			"pickup_lon":   ride.PickupLon,
			"dropoff_lat":  ride.DropoffLat,
			"dropoff_lon":  ride.DropoffLon,
			"status":       ride.Status,
		},
	}

	// Idempotency key prevents duplicate events
	idempotencyKey := "ride-requested-" + ride.ID

	// RULE 2: PublishIdempotent uses platform event-bus
	if err := p.eventBus.PublishIdempotent(ctx, idempotencyKey, event); err != nil {
		p.logger.Error("failed to publish ride_requested (idempotent)",
			zap.String("rideID", ride.ID),
			zap.Error(err))
		return err
	}

	p.logger.Info("ride_requested published (idempotent)",
		zap.String("rideID", ride.ID))

	return nil
}

// ===== EVENT CONSUMER PATTERNS =====

// NOTE: Event CONSUMPTION is handled by event-handler services
// Each service that needs to REACT to ride events implements:
// - packages/event-bus Subscriber
// - Consumes from shared/contracts/events topics
// - Examples: dispatch-service listens to ride.requested

// Example consumer (in dispatch-service or other service):
/*
type RideRequestedConsumer struct {
	handler event_bus.EventHandler
	logger  *zap.Logger
}

func (c *RideRequestedConsumer) Handle(ctx context.Context, evt *events.Event) error {
	// Extract data from event
	rideID := evt.Data["ride_id"].(string)
	passengerID := evt.Data["passenger_id"].(string)
	// ... call dispatch logic
	c.logger.Info("handled ride_requested", zap.String("rideID", rideID))
	return nil
}

// In bootstrap:
// eventBus.Subscribe(ctx, "ride.requested", consumer.Handle, event_bus.RetryPolicy{...})
*/
