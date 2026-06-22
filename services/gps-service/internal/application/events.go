// services/gps-service/internal/application/events.go
// GPS Service Event Publishing

package application

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
)

// GPSEventPublisher publishes GPS service events
type GPSEventPublisher struct {
	eventBus event_bus.EventBus
	logger   *zap.Logger
}

func NewGPSEventPublisher(eventBus event_bus.EventBus, logger *zap.Logger) *GPSEventPublisher {
	return &GPSEventPublisher{
		eventBus: eventBus,
		logger:   logger,
	}
}

// PublishDriverLocationUpdated publishes location update event
func (p *GPSEventPublisher) PublishDriverLocationUpdated(
	ctx context.Context,
	driverID string,
	latitude float64,
	longitude float64,
	accuracy float32,
) error {
	payload := map[string]interface{}{
		"driver_id":  driverID,
		"latitude":   latitude,
		"longitude":  longitude,
		"accuracy":   accuracy,
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "DriverLocationUpdated",
		AggregateID: driverID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published DriverLocationUpdated event",
		zap.String("driver_id", driverID),
		zap.Float64("latitude", latitude),
		zap.Float64("longitude", longitude))

	return nil
}

// PublishDriverOnline publishes driver online event
func (p *GPSEventPublisher) PublishDriverOnline(ctx context.Context, driverID string) error {
	payload := map[string]interface{}{
		"driver_id":  driverID,
		"timestamp":  time.Now().UTC(),
		"status":     "ONLINE",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "DriverOnline",
		AggregateID: driverID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published DriverOnline event", zap.String("driver_id", driverID))
	return nil
}

// PublishDriverOffline publishes driver offline event
func (p *GPSEventPublisher) PublishDriverOffline(ctx context.Context, driverID string) error {
	payload := map[string]interface{}{
		"driver_id":  driverID,
		"timestamp":  time.Now().UTC(),
		"status":     "OFFLINE",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "DriverOffline",
		AggregateID: driverID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published DriverOffline event", zap.String("driver_id", driverID))
	return nil
}

// PublishTripStarted publishes trip started event
func (p *GPSEventPublisher) PublishTripStarted(ctx context.Context, rideID, driverID string) error {
	payload := map[string]interface{}{
		"ride_id":    rideID,
		"driver_id":  driverID,
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "TripStarted",
		AggregateID: rideID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published TripStarted event",
		zap.String("ride_id", rideID),
		zap.String("driver_id", driverID))

	return nil
}

// PublishTripCompleted publishes trip completed event
func (p *GPSEventPublisher) PublishTripCompleted(
	ctx context.Context,
	rideID string,
	distance float32,
	duration int32,
) error {
	payload := map[string]interface{}{
		"ride_id":    rideID,
		"distance":   distance,
		"duration":   duration,
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "TripCompleted",
		AggregateID: rideID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published TripCompleted event",
		zap.String("ride_id", rideID),
		zap.Float32("distance", distance))

	return nil
}
