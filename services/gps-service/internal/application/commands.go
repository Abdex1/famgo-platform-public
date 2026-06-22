package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	"github.com/Abdex1/FamGo-platform/packages/telemetry"
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
)

// UpdateDriverLocationCommand represents updating a driver's location
type UpdateDriverLocationCommand struct {
	DriverID  string  `validate:"required,uuid"`
	Latitude  float64 `validate:"required,min=-90,max=90"`
	Longitude float64 `validate:"required,min=-180,max=180"`
	Accuracy  float32 `validate:"required,min=0"`
}

// UpdateDriverLocationHandler handles updating driver location
type UpdateDriverLocationHandler struct {
	locationRepo LocationRepository
	geofenceRepo GeofenceRepository
	eventBus     event_bus.EventBus
	service      *domain.LocationService
	metrics      telemetry.Metrics
	logger       telemetry.Logger
}

// NewUpdateDriverLocationHandler creates a new handler
func NewUpdateDriverLocationHandler(
	locationRepo LocationRepository,
	geofenceRepo GeofenceRepository,
	eventBus event_bus.EventBus,
	service *domain.LocationService,
	metrics telemetry.Metrics,
	logger telemetry.Logger,
) *UpdateDriverLocationHandler {
	return &UpdateDriverLocationHandler{
		locationRepo: locationRepo,
		geofenceRepo: geofenceRepo,
		eventBus:     eventBus,
		service:      service,
		metrics:      metrics,
		logger:       logger,
	}
}

// Handle executes the update location command
func (h *UpdateDriverLocationHandler) Handle(ctx context.Context, cmd UpdateDriverLocationCommand) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		h.metrics.RecordLatency("UpdateDriverLocation", duration)
	}()

	// 1. Get old location (if exists)
	oldLocation, _ := h.locationRepo.GetDriverLocation(ctx, cmd.DriverID)

	// 2. Create new location entity
	newLocation := domain.NewDriverLocation(cmd.DriverID, cmd.Latitude, cmd.Longitude, cmd.Accuracy)

	// 3. Check for geofence changes (domain logic)
	geofences, err := h.geofenceRepo.GetAllGeofences(ctx)
	if err != nil {
		h.logger.Error("Failed to get geofences", map[string]interface{}{
			"driver_id": cmd.DriverID,
			"error":     err.Error(),
		})
		h.metrics.RecordError("UpdateDriverLocation", err)
		return err
	}

	// 4. Check geofence entries/exits
	for _, geofence := range geofences {
		wasInside := oldLocation != nil && h.service.IsWithinGeofence(*oldLocation, geofence)
		isInside := h.service.IsWithinGeofence(*newLocation, geofence)

		// Geofence entered
		if !wasInside && isInside {
			_ = h.eventBus.PublishIdempotent(ctx, events.GeofenceEnteredEvent{
				EventID:     uuid.New().String(),
				EventType:   events.EventTypeGeofenceEntered,
				Version:     events.VersionV1,
				AggregateID: cmd.DriverID,
				Timestamp:   time.Now(),
				Data: map[string]interface{}{
					"geofence_id": geofence.ID,
					"latitude":    newLocation.Latitude,
					"longitude":   newLocation.Longitude,
				},
			})
		}

		// Geofence exited
		if wasInside && !isInside {
			_ = h.eventBus.PublishIdempotent(ctx, events.GeofenceExitedEvent{
				EventID:     uuid.New().String(),
				EventType:   events.EventTypeGeofenceExited,
				Version:     events.VersionV1,
				AggregateID: cmd.DriverID,
				Timestamp:   time.Now(),
				Data: map[string]interface{}{
					"geofence_id": geofence.ID,
					"latitude":    newLocation.Latitude,
					"longitude":   newLocation.Longitude,
				},
			})
		}
	}

	// 5. Persist new location
	if err := h.locationRepo.UpdateDriverLocation(ctx, newLocation); err != nil {
		h.logger.Error("Failed to update location", map[string]interface{}{
			"driver_id": cmd.DriverID,
			"error":     err.Error(),
		})
		h.metrics.RecordError("UpdateDriverLocation", err)
		return err
	}

	// 6. Publish event through shared/contracts/events
	err = h.eventBus.PublishIdempotent(ctx, events.DriverLocationUpdatedEvent{
		EventID:     uuid.New().String(),
		EventType:   events.EventTypeDriverLocationUpdated,
		Version:     events.VersionV1,
		AggregateID: cmd.DriverID,
		Timestamp:   time.Now(),
		Data: map[string]interface{}{
			"latitude":  newLocation.Latitude,
			"longitude": newLocation.Longitude,
			"accuracy":  newLocation.Accuracy,
		},
	})

	if err != nil {
		h.logger.Error("Failed to publish event", map[string]interface{}{
			"driver_id": cmd.DriverID,
			"error":     err.Error(),
		})
		h.metrics.RecordError("PublishLocationEvent", err)
		return err
	}

	h.metrics.RecordSuccess("UpdateDriverLocation")
	h.logger.Info("Location updated", map[string]interface{}{
		"driver_id": cmd.DriverID,
		"latitude":  newLocation.Latitude,
		"longitude": newLocation.Longitude,
	})

	return nil
}
