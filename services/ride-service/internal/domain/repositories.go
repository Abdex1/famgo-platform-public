// services/ride-service/internal/domain/repositories.go
// Repository interfaces

package domain

import "context"

// RideRepository defines ride persistence operations
type RideRepository interface {
	GetRide(ctx context.Context, rideID string) (*Ride, error)
	CreateRide(ctx context.Context, ride *Ride) error
	UpdateRide(ctx context.Context, ride *Ride) error
	GetRidesByPassenger(ctx context.Context, passengerID string, limit, offset int) ([]Ride, error)
	GetRidesByDriver(ctx context.Context, driverID string, limit, offset int) ([]Ride, error)
	GetActiveRides(ctx context.Context) ([]Ride, error)
}

// RideStatusHistoryRepository tracks state transitions
type RideStatusHistoryRepository interface {
	LogTransition(ctx context.Context, history *RideStatusHistory) error
	GetHistory(ctx context.Context, rideID string) ([]RideStatusHistory, error)
}

// RideCache defines caching operations
type RideCache interface {
	GetRide(ctx context.Context, rideID string) (*Ride, error)
	SetRide(ctx context.Context, ride *Ride, ttl int32) error
	DeleteRide(ctx context.Context, rideID string) error
	GetActiveRides(ctx context.Context) ([]Ride, error)
	SetActiveRides(ctx context.Context, rides []Ride, ttl int32) error
}
