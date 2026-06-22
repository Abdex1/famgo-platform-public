package application

import (
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
)

// LocationRepository interface (what application depends on)
type LocationRepository interface {
	GetDriverLocation(ctx interface{}, driverID string) (*domain.DriverLocation, error)
	UpdateDriverLocation(ctx interface{}, location *domain.DriverLocation) error
	ListActiveLocations(ctx interface{}) ([]domain.DriverLocation, error)
	DeleteDriverLocation(ctx interface{}, driverID string) error
}

// TripRepository interface (what application depends on)
type TripRepository interface {
	GetTrip(ctx interface{}, tripID string) (*domain.Trip, error)
	CreateTrip(ctx interface{}, trip *domain.Trip) error
	UpdateTrip(ctx interface{}, trip *domain.Trip) error
	AddRoutePoint(ctx interface{}, tripID string, point domain.RoutePoint) error
	GetTripsByDriver(ctx interface{}, driverID string) ([]domain.Trip, error)
}

// GeofenceRepository interface (what application depends on)
type GeofenceRepository interface {
	GetGeofence(ctx interface{}, geofenceID string) (*domain.Geofence, error)
	GetAllGeofences(ctx interface{}) ([]domain.Geofence, error)
	CreateGeofence(ctx interface{}, geofence *domain.Geofence) error
	GetGeofencesByPoint(ctx interface{}, latitude float64, longitude float64) ([]domain.Geofence, error)
}
