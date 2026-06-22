// services/ride-service/internal/application/interfaces.go
// Application Layer Interfaces

package application

import (
	"context"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// RideRepository interface
type RideRepository interface {
	GetRide(ctx context.Context, rideID string) (*domain.Ride, error)
	CreateRide(ctx context.Context, ride *domain.Ride) error
	UpdateRide(ctx context.Context, ride *domain.Ride) error
	GetRidesByPassenger(ctx context.Context, passengerID string, limit, offset int) ([]domain.Ride, error)
	GetRidesByDriver(ctx context.Context, driverID string, limit, offset int) ([]domain.Ride, error)
	GetActiveRides(ctx context.Context) ([]domain.Ride, error)
}

// RideCache interface
type RideCache interface {
	GetRide(ctx context.Context, rideID string) (*domain.Ride, error)
	SetRide(ctx context.Context, ride *domain.Ride, ttl int32) error
	DeleteRide(ctx context.Context, rideID string) error
	GetActiveRides(ctx context.Context) ([]domain.Ride, error)
	SetActiveRides(ctx context.Context, rides []domain.Ride, ttl int32) error
}
