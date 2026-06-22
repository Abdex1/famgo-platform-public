package infrastructure

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	"github.com/Abdex1/FamGo-platform/packages/telemetry"
)

// EventPublisher handles publishing GPS events through the platform event-bus
type EventPublisher struct {
	eventBus event_bus.EventBus
	metrics  telemetry.Metrics
	logger   telemetry.Logger
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(
	eventBus event_bus.EventBus,
	metrics telemetry.Metrics,
	logger telemetry.Logger,
) *EventPublisher {
	return &EventPublisher{
		eventBus: eventBus,
		metrics:  metrics,
		logger:   logger,
	}
}

// PublishLocationUpdated publishes location update event
func (p *EventPublisher) PublishLocationUpdated(ctx context.Context, event interface{}) error {
	err := p.eventBus.PublishIdempotent(ctx, event)
	if err != nil {
		p.metrics.RecordError("PublishLocationUpdated", err)
		p.logger.Error("Failed to publish location updated event", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to publish location updated event: %w", err)
	}

	p.metrics.RecordSuccess("PublishLocationUpdated")
	return nil
}

// PublishGeofenceEntered publishes geofence entry event
func (p *EventPublisher) PublishGeofenceEntered(ctx context.Context, event interface{}) error {
	err := p.eventBus.PublishIdempotent(ctx, event)
	if err != nil {
		p.metrics.RecordError("PublishGeofenceEntered", err)
		p.logger.Error("Failed to publish geofence entered event", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to publish geofence entered event: %w", err)
	}

	p.metrics.RecordSuccess("PublishGeofenceEntered")
	return nil
}

// PublishGeofenceExited publishes geofence exit event
func (p *EventPublisher) PublishGeofenceExited(ctx context.Context, event interface{}) error {
	err := p.eventBus.PublishIdempotent(ctx, event)
	if err != nil {
		p.metrics.RecordError("PublishGeofenceExited", err)
		p.logger.Error("Failed to publish geofence exited event", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to publish geofence exited event: %w", err)
	}

	p.metrics.RecordSuccess("PublishGeofenceExited")
	return nil
}

// PublishTripStarted publishes trip started event
func (p *EventPublisher) PublishTripStarted(ctx context.Context, event interface{}) error {
	err := p.eventBus.PublishIdempotent(ctx, event)
	if err != nil {
		p.metrics.RecordError("PublishTripStarted", err)
		p.logger.Error("Failed to publish trip started event", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to publish trip started event: %w", err)
	}

	p.metrics.RecordSuccess("PublishTripStarted")
	return nil
}

// PublishTripCompleted publishes trip completed event
func (p *EventPublisher) PublishTripCompleted(ctx context.Context, event interface{}) error {
	err := p.eventBus.PublishIdempotent(ctx, event)
	if err != nil {
		p.metrics.RecordError("PublishTripCompleted", err)
		p.logger.Error("Failed to publish trip completed event", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to publish trip completed event: %w", err)
	}

	p.metrics.RecordSuccess("PublishTripCompleted")
	return nil
}
