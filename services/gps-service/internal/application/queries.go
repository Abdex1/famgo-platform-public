package application

import (
	"context"
	"time"

	"github.com/Abdex1/FamGo-platform/packages/telemetry"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
)

// GetDriverLocationQuery retrieves current driver location
type GetDriverLocationQuery struct {
	DriverID string `validate:"required,uuid"`
}

// GetDriverLocationHandler handles location retrieval
type GetDriverLocationHandler struct {
	locationRepo LocationRepository
	metrics      telemetry.Metrics
	logger       telemetry.Logger
}

// NewGetDriverLocationHandler creates a new handler
func NewGetDriverLocationHandler(
	locationRepo LocationRepository,
	metrics telemetry.Metrics,
	logger telemetry.Logger,
) *GetDriverLocationHandler {
	return &GetDriverLocationHandler{
		locationRepo: locationRepo,
		metrics:      metrics,
		logger:       logger,
	}
}

// Handle executes the get location query
func (h *GetDriverLocationHandler) Handle(ctx context.Context, q GetDriverLocationQuery) (*domain.DriverLocation, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		h.metrics.RecordLatency("GetDriverLocation", duration)
	}()

	location, err := h.locationRepo.GetDriverLocation(ctx, q.DriverID)
	if err != nil {
		h.metrics.RecordError("GetDriverLocation", err)
		h.logger.Error("Failed to get location", map[string]interface{}{
			"driver_id": q.DriverID,
			"error":     err.Error(),
		})
		return nil, err
	}

	h.metrics.RecordSuccess("GetDriverLocation")
	return location, nil
}

// ========================================

// GetNearbyDriversQuery retrieves drivers near a location
type GetNearbyDriversQuery struct {
	Latitude  float64 `validate:"required,min=-90,max=90"`
	Longitude float64 `validate:"required,min=-180,max=180"`
	RadiusM   float64 `validate:"required,min=100,max=50000"` // 100m to 50km
}

// GetNearbyDriversHandler handles nearby driver retrieval
type GetNearbyDriversHandler struct {
	locationRepo LocationRepository
	service      *domain.LocationService
	metrics      telemetry.Metrics
	logger       telemetry.Logger
}

// NewGetNearbyDriversHandler creates a new handler
func NewGetNearbyDriversHandler(
	locationRepo LocationRepository,
	service *domain.LocationService,
	metrics telemetry.Metrics,
	logger telemetry.Logger,
) *GetNearbyDriversHandler {
	return &GetNearbyDriversHandler{
		locationRepo: locationRepo,
		service:      service,
		metrics:      metrics,
		logger:       logger,
	}
}

// Handle executes the get nearby drivers query
func (h *GetNearbyDriversHandler) Handle(ctx context.Context, q GetNearbyDriversQuery) ([]domain.DriverLocation, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		h.metrics.RecordLatency("GetNearbyDrivers", duration)
	}()

	// Get all active locations
	allLocations, err := h.locationRepo.ListActiveLocations(ctx)
	if err != nil {
		h.metrics.RecordError("GetNearbyDrivers", err)
		h.logger.Error("Failed to get locations", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Filter by radius (using domain logic)
	var nearby []domain.DriverLocation
	for _, loc := range allLocations {
		distance := h.service.CalculateDistance(
			q.Latitude, q.Longitude,
			loc.Latitude, loc.Longitude,
		)

		if distance <= q.RadiusM {
			nearby = append(nearby, loc)
		}
	}

	h.metrics.RecordSuccess("GetNearbyDrivers")
	h.logger.Info("Found nearby drivers", map[string]interface{}{
		"count":  len(nearby),
		"radius": q.RadiusM,
	})

	return nearby, nil
}
