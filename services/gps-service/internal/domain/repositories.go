package domain

import (
	"context"
)

// LocationRepository defines what the application layer needs from location storage
type LocationRepository interface {
	// GetDriverLocation retrieves current location for a driver
	GetDriverLocation(ctx context.Context, driverID string) (*DriverLocation, error)

	// UpdateDriverLocation updates or creates driver location
	UpdateDriverLocation(ctx context.Context, location *DriverLocation) error

	// ListActiveLocations retrieves all active driver locations
	ListActiveLocations(ctx context.Context) ([]DriverLocation, error)

	// DeleteDriverLocation removes driver location (e.g., on logout)
	DeleteDriverLocation(ctx context.Context, driverID string) error
}

// TripRepository defines what the application layer needs from trip storage
type TripRepository interface {
	// GetTrip retrieves a trip by ID
	GetTrip(ctx context.Context, tripID string) (*Trip, error)

	// CreateTrip creates a new trip
	CreateTrip(ctx context.Context, trip *Trip) error

	// UpdateTrip updates an existing trip
	UpdateTrip(ctx context.Context, trip *Trip) error

	// AddRoutePoint adds a location point to a trip's route
	AddRoutePoint(ctx context.Context, tripID string, point RoutePoint) error

	// GetTripsByDriver retrieves all active trips for a driver
	GetTripsByDriver(ctx context.Context, driverID string) ([]Trip, error)
}

// GeofenceRepository defines what the application layer needs from geofence storage
type GeofenceRepository interface {
	// GetGeofence retrieves a geofence by ID
	GetGeofence(ctx context.Context, geofenceID string) (*Geofence, error)

	// GetAllGeofences retrieves all geofences
	GetAllGeofences(ctx context.Context) ([]Geofence, error)

	// CreateGeofence creates a new geofence
	CreateGeofence(ctx context.Context, geofence *Geofence) error

	// GetGeofencesByPoint retrieves geofences that contain a point
	GetGeofencesByPoint(ctx context.Context, latitude float64, longitude float64) ([]Geofence, error)
}
