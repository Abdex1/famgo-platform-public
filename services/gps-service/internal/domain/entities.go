package domain

import (
	"time"
)

// DriverLocation represents the current location of a driver
type DriverLocation struct {
	ID        string    // UUID
	DriverID  string    // Foreign key to driver
	Latitude  float64   // Geographic latitude
	Longitude float64   // Geographic longitude
	Accuracy  float32   // Accuracy in meters
	UpdatedAt time.Time
}

// Trip represents an active trip being tracked
type Trip struct {
	ID        string       // UUID
	RideID    string       // Foreign key to ride
	DriverID  string       // Which driver
	StartedAt time.Time
	Location  DriverLocation
	Route     []RoutePoint
	Status    TripStatus
}

// TripStatus represents the state of a trip
type TripStatus string

const (
	TripStatusActive    TripStatus = "ACTIVE"
	TripStatusCompleted TripStatus = "COMPLETED"
	TripStatusCancelled TripStatus = "CANCELLED"
)

// Geofence represents a geographic boundary (zone, area)
type Geofence struct {
	ID        string
	Name      string
	Latitude  float64 // Center point
	Longitude float64 // Center point
	Radius    float32 // Radius in meters
	CreatedAt time.Time
}

// RoutePoint represents a single point in a trip route
type RoutePoint struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// NewDriverLocationWithID creates a new driver location entity
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewDriverLocationWithID(id, driverID string, lat float64, lon float64, accuracy float32) *DriverLocation {
	return &DriverLocation{
		ID:        id,
		DriverID:  driverID,
		Latitude:  lat,
		Longitude: lon,
		Accuracy:  accuracy,
		UpdatedAt: time.Now(),
	}
}

// NewTripWithID creates a new trip entity
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewTripWithID(id, rideID, driverID string) *Trip {
	return &Trip{
		ID:        id,
		RideID:    rideID,
		DriverID:  driverID,
		StartedAt: time.Now(),
		Status:    TripStatusActive,
		Route:     []RoutePoint{},
	}
}

// NewGeofenceWithID creates a new geofence entity
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewGeofenceWithID(id, name string, lat float64, lon float64, radius float32) *Geofence {
	return &Geofence{
		ID:        id,
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
		Radius:    radius,
		CreatedAt: time.Now(),
	}
}

// AddRoutePoint adds a point to the trip route
func (t *Trip) AddRoutePoint(lat float64, lon float64) {
	t.Route = append(t.Route, RoutePoint{
		Latitude:  lat,
		Longitude: lon,
		Timestamp: time.Now(),
	})
}

// Complete marks trip as completed
func (t *Trip) Complete() {
	t.Status = TripStatusCompleted
}

// Cancel marks trip as cancelled
func (t *Trip) Cancel() {
	t.Status = TripStatusCancelled
}
