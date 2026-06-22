// services/gps-service/internal/application/usecases/location_usecases.go
// Application use cases for location management

package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/redis"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/repositories"
)

// LocationUseCases holds all location-related use cases
type LocationUseCases struct {
	locationRepository   *repositories.DriverLocationRepository
	geoIndexStore        *redis.GeoIndexStore
	driverTrackingStore  *redis.DriverTrackingStore
	locationService      *services.LocationService
	redisGeoService      *services.RedisGeoService
}

// NewLocationUseCases creates new location use cases
func NewLocationUseCases(
	locationRepo *repositories.DriverLocationRepository,
	geoIndexStore *redis.GeoIndexStore,
	trackingStore *redis.DriverTrackingStore,
	locationService *services.LocationService,
	redisGeoService *services.RedisGeoService,
) *LocationUseCases {
	return &LocationUseCases{
		locationRepository:  locationRepo,
		geoIndexStore:       geoIndexStore,
		driverTrackingStore: trackingStore,
		locationService:     locationService,
		redisGeoService:     redisGeoService,
	}
}

// UpdateLocationInput holds input for updating driver location
type UpdateLocationInput struct {
	DriverID  string
	Latitude  float64
	Longitude float64
	Altitude  float64
	Accuracy  float64
	Speed     float64
	Heading   float64
	Timestamp int64
}

// UpdateLocationOutput holds output after location update
type UpdateLocationOutput struct {
	Success           bool
	DriverID          string
	Distance          float64
	ETAMinutes        float64
	QualityIssues     string
	ConsecutiveFails  int
}

// UpdateDriverLocation updates driver location and propagates to caches
func (uc *LocationUseCases) UpdateDriverLocation(
	ctx context.Context,
	input *UpdateLocationInput,
) (*UpdateLocationOutput, error) {
	if input == nil || input.DriverID == "" {
		return nil, fmt.Errorf("invalid update location input")
	}

	// Create new location value object
	newLocation, err := valueobjects.NewGeolocation(
		input.Latitude,
		input.Longitude,
		input.Altitude,
		input.Accuracy,
		input.Speed,
		input.Heading,
		input.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid location data: %w", err)
	}

	// Get existing driver location
	driverLocation, err := uc.locationRepository.GetByDriverID(ctx, input.DriverID)
	if err != nil {
		// If not exists, create new
		driverLocation, err = entities.NewDriverLocation(input.DriverID, newLocation)
		if err != nil {
			return nil, fmt.Errorf("failed to create driver location: %w", err)
		}
		if err := uc.locationRepository.Create(ctx, driverLocation); err != nil {
			return nil, fmt.Errorf("failed to persist new driver location: %w", err)
		}
	} else {
		// Update existing location
		previousLocation := driverLocation.CurrentLocation
		if err := driverLocation.UpdateLocation(newLocation); err != nil {
			return nil, fmt.Errorf("failed to update driver location: %w", err)
		}

		// Analyze trajectory for anomalies
		if previousLocation != nil {
			trajectory := uc.locationService.AnalyzeTrajectory(previousLocation, newLocation, 200) // 200 km/h max
			if !trajectory.IsValid || trajectory.IsAnomaly {
				driverLocation.RecordFailure()
				output := &UpdateLocationOutput{
					Success:          false,
					DriverID:         input.DriverID,
					QualityIssues:    fmt.Sprintf("Anomaly detected: %.2f km/h", trajectory.CalculatedSpeed),
					ConsecutiveFails: driverLocation.ConsecutiveFailures,
				}
				return output, nil
			}
		}

		// Reset failures on successful update
		driverLocation.ResetFailures()

		// Update location hash
		hash := fmt.Sprintf("%d", time.Now().UnixNano())
		driverLocation.UpdateLocationHash(hash)

		// Update geohash
		geohash, _ := uc.locationService.CalculateGeohash(input.Latitude, input.Longitude, 8)
		driverLocation.UpdateGeohash(geohash)

		if err := uc.locationRepository.Update(ctx, driverLocation); err != nil {
			return nil, fmt.Errorf("failed to persist location update: %w", err)
		}
	}

	// Update GEO index
	if err := uc.geoIndexStore.AddLocation(ctx, input.DriverID, newLocation); err != nil {
		// Log error but don't fail
		fmt.Printf("warning: failed to update GEO index: %v\n", err)
	}

	// Update tracking store
	if err := uc.driverTrackingStore.UpdateDriverLocation(ctx, input.DriverID, newLocation); err != nil {
		// Log error but don't fail
		fmt.Printf("warning: failed to update tracking store: %v\n", err)
	}

	output := &UpdateLocationOutput{
		Success:           true,
		DriverID:          input.DriverID,
		ConsecutiveFails:  driverLocation.ConsecutiveFailures,
	}

	// Calculate distance from previous location if available
	if driverLocation.PreviousLocation != nil {
		output.Distance = driverLocation.CurrentLocation.DistanceToKm(driverLocation.PreviousLocation)
	}

	return output, nil
}

// FindNearbyDriversInput holds input for finding nearby drivers
type FindNearbyDriversInput struct {
	Latitude      float64
	Longitude     float64
	RadiusKm      float64
	Limit         int
	BaseSpeedKmH  float64
	OnlyOnline    bool
}

// NearbyDriverResult represents a nearby driver
type NearbyDriverResult struct {
	DriverID          string
	Distance          float64
	ETAMinutes        float64
	Bearing           float64
	Latitude          float64
	Longitude         float64
	IsOnline          bool
	AcceptanceRate    float64
	Rating            float64
	RidesCompleted    int
}

// FindNearbyDrivers finds nearby drivers using GEO index
func (uc *LocationUseCases) FindNearbyDrivers(
	ctx context.Context,
	input *FindNearbyDriversInput,
) ([]NearbyDriverResult, error) {
	if input == nil || input.RadiusKm <= 0 {
		return nil, fmt.Errorf("invalid find nearby drivers input")
	}

	referenceLocation, err := valueobjects.NewGeolocation(
		input.Latitude,
		input.Longitude,
		0, 0, 0, 0,
		time.Now().UnixMilli(),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid reference location: %w", err)
	}

	// Query GEO index
	driverIDs, err := uc.geoIndexStore.FindNearby(
		ctx,
		referenceLocation,
		input.RadiusKm,
		input.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby drivers: %w", err)
	}

	if len(driverIDs) == 0 {
		return []NearbyDriverResult{}, nil
	}

	var results []NearbyDriverResult

	for _, driverID := range driverIDs {
		driverLocation, err := uc.locationRepository.GetByDriverID(ctx, driverID)
		if err != nil {
			continue // Skip drivers that can't be retrieved
		}

		if input.OnlyOnline && !driverLocation.IsActive() {
			continue
		}

		distance := referenceLocation.DistanceToKm(driverLocation.CurrentLocation)
		eta := referenceLocation.EstimatedArrivalTime(driverLocation.CurrentLocation, input.BaseSpeedKmH)
		bearing := referenceLocation.BearingTo(driverLocation.CurrentLocation)

		// Get tracking info
		trackingStatus, _ := uc.driverTrackingStore.GetDriverStatus(ctx, driverID)

		result := NearbyDriverResult{
			DriverID:       driverID,
			Distance:       distance,
			ETAMinutes:     eta,
			Bearing:        bearing,
			Latitude:       driverLocation.CurrentLocation.Coordinates.Latitude,
			Longitude:      driverLocation.CurrentLocation.Coordinates.Longitude,
			IsOnline:       driverLocation.IsActive(),
			Rating:         driverLocation.AverageRating,
			RidesCompleted: driverLocation.CompletedRideCount,
		}

		if trackingStatus != nil {
			result.AcceptanceRate = trackingStatus.AcceptanceRate
		}

		results = append(results, result)
	}

	// Sort by distance
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Distance < results[i].Distance {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

// GetDriverLocationOutput holds output for getting driver location
type GetDriverLocationOutput struct {
	DriverID         string
	Latitude         float64
	Longitude        float64
	Altitude         float64
	Accuracy         float64
	Speed            float64
	Heading          float64
	Status           string
	IsOnline         bool
	LastUpdateAt     time.Time
	AcceptanceRate   float64
	Rating           float64
	RidesCompleted   int
}

// GetDriverLocation retrieves current driver location
func (uc *LocationUseCases) GetDriverLocation(
	ctx context.Context,
	driverID string,
) (*GetDriverLocationOutput, error) {
	if driverID == "" {
		return nil, fmt.Errorf("invalid driver ID")
	}

	driverLocation, err := uc.locationRepository.GetByDriverID(ctx, driverID)
	if err != nil {
		return nil, fmt.Errorf("driver location not found: %w", err)
	}

	trackingStatus, _ := uc.driverTrackingStore.GetDriverStatus(ctx, driverID)

	output := &GetDriverLocationOutput{
		DriverID:       driverID,
		Latitude:       driverLocation.CurrentLocation.Coordinates.Latitude,
		Longitude:      driverLocation.CurrentLocation.Coordinates.Longitude,
		Altitude:       driverLocation.CurrentLocation.Altitude,
		Accuracy:       driverLocation.CurrentLocation.Accuracy,
		Speed:          driverLocation.CurrentLocation.Speed,
		Heading:        driverLocation.CurrentLocation.Heading,
		Status:         string(driverLocation.Status),
		IsOnline:       driverLocation.IsActive(),
		LastUpdateAt:   driverLocation.LastUpdateAt,
		Rating:         driverLocation.AverageRating,
		RidesCompleted: driverLocation.CompletedRideCount,
	}

	if trackingStatus != nil {
		output.AcceptanceRate = trackingStatus.AcceptanceRate
	}

	return output, nil
}

// UpdateDriverStatusInput holds input for updating driver status
type UpdateDriverStatusInput struct {
	DriverID string
	Status   string // "online", "offline", "on_ride", "break", "maintenance"
}

// UpdateDriverStatus updates driver status
func (uc *LocationUseCases) UpdateDriverStatus(
	ctx context.Context,
	input *UpdateDriverStatusInput,
) error {
	if input == nil || input.DriverID == "" {
		return fmt.Errorf("invalid update driver status input")
	}

	driverLocation, err := uc.locationRepository.GetByDriverID(ctx, input.DriverID)
	if err != nil {
		return fmt.Errorf("driver location not found: %w", err)
	}

	// Update status
	if err := driverLocation.SetStatus(entities.DriverStatus(input.Status)); err != nil {
		return fmt.Errorf("failed to set driver status: %w", err)
	}

	// Update online flag
	if input.Status == "online" {
		driverLocation.GoOnline()
		if err := uc.driverTrackingStore.SetDriverOnline(ctx, input.DriverID, driverLocation.CurrentLocation); err != nil {
			fmt.Printf("warning: failed to update tracking store: %v\n", err)
		}
	} else {
		driverLocation.GoOffline()
		if err := uc.driverTrackingStore.SetDriverOffline(ctx, input.DriverID); err != nil {
			fmt.Printf("warning: failed to update tracking store: %v\n", err)
		}
	}

	// Persist
	if err := uc.locationRepository.Update(ctx, driverLocation); err != nil {
		return fmt.Errorf("failed to update driver status: %w", err)
	}

	return nil
}

// BulkUpdateLocationsInput holds input for bulk updating locations
type BulkUpdateLocationsInput struct {
	Locations []UpdateLocationInput
}

// BulkUpdateLocationsOutput holds output after bulk updates
type BulkUpdateLocationsOutput struct {
	Processed   int
	Failed      int
	FailedDrivers []string
}

// BulkUpdateLocations updates multiple driver locations efficiently
func (uc *LocationUseCases) BulkUpdateLocations(
	ctx context.Context,
	input *BulkUpdateLocationsInput,
) (*BulkUpdateLocationsOutput, error) {
	if input == nil || len(input.Locations) == 0 {
		return nil, fmt.Errorf("invalid bulk update input")
	}

	output := &BulkUpdateLocationsOutput{}
	locations := make(map[string]*valueobjects.Geolocation)

	for _, locInput := range input.Locations {
		newLocation, err := valueobjects.NewGeolocation(
			locInput.Latitude,
			locInput.Longitude,
			locInput.Altitude,
			locInput.Accuracy,
			locInput.Speed,
			locInput.Heading,
			locInput.Timestamp,
		)
		if err != nil {
			output.Failed++
			output.FailedDrivers = append(output.FailedDrivers, locInput.DriverID)
			continue
		}

		locations[locInput.DriverID] = newLocation
		output.Processed++
	}

	// Batch update GEO index
	if len(locations) > 0 {
		if err := uc.geoIndexStore.AddLocations(ctx, locations); err != nil {
			fmt.Printf("warning: failed to batch update GEO index: %v\n", err)
		}
	}

	return output, nil
}
