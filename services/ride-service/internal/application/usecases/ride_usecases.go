// services/ride-service/internal/application/usecases/ride_usecases.go
package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/infrastructure/repositories"
)

type RideUseCases struct {
	repo        *repositories.RideRepository
	rideService *services.RideService
}

func NewRideUseCases(repo *repositories.RideRepository, svc *services.RideService) *RideUseCases {
	return &RideUseCases{repo: repo, rideService: svc}
}

type CreateRideInput struct {
	RiderID          string
	PickupLat        float64
	PickupLng        float64
	DropoffLat       float64
	DropoffLng       float64
	PickupAddress    string
	DropoffAddress   string
	PaymentMethod    string
	PassengerCount   int
	RideType         string
}

type CreateRideOutput struct {
	RideID            string
	EstimatedDistance float64
	EstimatedDuration time.Duration
	EstimatedFare     float64
}

func (uc *RideUseCases) CreateRideRequest(ctx context.Context, input *CreateRideInput) (*CreateRideOutput, error) {
	if input == nil || input.RiderID == "" {
		return nil, fmt.Errorf("invalid create ride input")
	}

	valid, errMsg := uc.rideService.ValidateRideRequest(
		input.PickupLat, input.PickupLng,
		input.DropoffLat, input.DropoffLng,
		0.5, 100.0,
	)
	if !valid {
		return nil, fmt.Errorf("invalid ride request: %s", errMsg)
	}

	ride, err := entities.NewRide(input.RiderID, input.PickupLat, input.PickupLng, input.DropoffLat, input.DropoffLng)
	if err != nil {
		return nil, fmt.Errorf("failed to create ride entity: %w", err)
	}

	ride.PickupAddress = input.PickupAddress
	ride.DropoffAddress = input.DropoffAddress
	ride.PaymentMethod = input.PaymentMethod
	ride.PassengerCount = input.PassengerCount
	ride.RideType = input.RideType

	if err := uc.repo.Create(ctx, ride); err != nil {
		return nil, fmt.Errorf("failed to persist ride: %w", err)
	}

	// Calculate fare
	fareCalc := uc.rideService.CalculateFare(ride.EstimatedDistance, ride.EstimatedDuration, 1.0, 0)
	ride.EstimatedFare = fareCalc.TotalFare

	return &CreateRideOutput{
		RideID:            ride.ID,
		EstimatedDistance: ride.EstimatedDistance,
		EstimatedDuration: ride.EstimatedDuration,
		EstimatedFare:     fareCalc.TotalFare,
	}, nil
}

type AcceptRideInput struct {
	RideID   string
	DriverID string
}

func (uc *RideUseCases) AcceptRide(ctx context.Context, input *AcceptRideInput) error {
	ride, err := uc.repo.GetByID(ctx, input.RideID)
	if err != nil {
		return fmt.Errorf("ride not found: %w", err)
	}

	if err := ride.Accept(input.DriverID); err != nil {
		return fmt.Errorf("cannot accept ride: %w", err)
	}

	return uc.repo.Update(ctx, ride)
}

type UpdateRideStatusInput struct {
	RideID string
	Status string
}

func (uc *RideUseCases) UpdateRideStatus(ctx context.Context, input *UpdateRideStatusInput) error {
	ride, err := uc.repo.GetByID(ctx, input.RideID)
	if err != nil {
		return fmt.Errorf("ride not found: %w", err)
	}

	targetStatus := entities.RideStatus(input.Status)

	// Validate transition
	if !uc.rideService.CanTransitionToStatus(ride.Status, targetStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", ride.Status, targetStatus)
	}

	switch targetStatus {
	case entities.StatusPickedUp:
		if err := ride.PickupPassenger(); err != nil {
			return fmt.Errorf("cannot pickup passenger: %w", err)
		}
	case entities.StatusOngoing:
		if err := ride.StartRide(); err != nil {
			return fmt.Errorf("cannot start ride: %w", err)
		}
	case entities.StatusCancelled:
		if err := ride.Cancel("driver_cancelled"); err != nil {
			return fmt.Errorf("cannot cancel ride: %w", err)
		}
	}

	return uc.repo.Update(ctx, ride)
}

type CompleteRideInput struct {
	RideID        string
	ActualDistance float64
	ActualFare    float64
	DriverRating  int
}

func (uc *RideUseCases) CompleteRide(ctx context.Context, input *CompleteRideInput) error {
	ride, err := uc.repo.GetByID(ctx, input.RideID)
	if err != nil {
		return fmt.Errorf("ride not found: %w", err)
	}

	if err := ride.CompleteRide(input.ActualDistance, input.ActualFare); err != nil {
		return fmt.Errorf("cannot complete ride: %w", err)
	}

	if input.DriverRating > 0 {
		if err := ride.SetDriverRating(input.DriverRating); err != nil {
			return fmt.Errorf("cannot set rating: %w", err)
		}
	}

	return uc.repo.Update(ctx, ride)
}

type GetRideDetailsOutput struct {
	RideID            string
	RiderID           string
	DriverID          *string
	Status            string
	EstimatedDistance float64
	EstimatedFare     float64
	ActualFare        *float64
	Distance          float64
	Duration          float64
	DriverRating      *int
}

func (uc *RideUseCases) GetRideDetails(ctx context.Context, rideID string) (*GetRideDetailsOutput, error) {
	ride, err := uc.repo.GetByID(ctx, rideID)
	if err != nil {
		return nil, fmt.Errorf("ride not found: %w", err)
	}

	metrics := uc.rideService.CalculateMetrics(ride)

	return &GetRideDetailsOutput{
		RideID:            ride.ID,
		RiderID:           ride.RiderID,
		DriverID:          ride.DriverID,
		Status:            string(ride.Status),
		EstimatedDistance: ride.EstimatedDistance,
		EstimatedFare:     ride.EstimatedFare,
		ActualFare:        ride.ActualFare,
		Distance:          metrics.Distance,
		Duration:          metrics.Duration.Minutes(),
		DriverRating:      ride.DriverRating,
	}, nil
}

type CancelRideInput struct {
	RideID string
	Reason string
}

func (uc *RideUseCases) CancelRide(ctx context.Context, input *CancelRideInput) error {
	ride, err := uc.repo.GetByID(ctx, input.RideID)
	if err != nil {
		return fmt.Errorf("ride not found: %w", err)
	}

	if err := ride.Cancel(input.Reason); err != nil {
		return fmt.Errorf("cannot cancel ride: %w", err)
	}

	return uc.repo.Update(ctx, ride)
}
