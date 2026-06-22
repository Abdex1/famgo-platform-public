// services/ride-service/internal/infrastructure/repositories/ride_repository.go
package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/entities"
)

type RideRepository struct {
	pool *pgxpool.Pool
}

func NewRideRepository(pool *pgxpool.Pool) *RideRepository {
	return &RideRepository{pool: pool}
}

func (r *RideRepository) Create(ctx context.Context, ride *entities.Ride) error {
	query := `
		INSERT INTO rides (
			id, rider_id, pickup_latitude, pickup_longitude, dropoff_latitude, dropoff_longitude,
			estimated_distance, estimated_duration, status, requested_at, ride_type, passenger_count,
			estimated_fare, payment_method, payment_status, pickup_address, dropoff_address,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	err := r.pool.QueryRow(ctx, query,
		ride.ID, ride.RiderID,
		ride.PickupLatitude, ride.PickupLongitude,
		ride.DropoffLatitude, ride.DropoffLongitude,
		ride.EstimatedDistance, ride.EstimatedDuration.Minutes(),
		string(ride.Status), ride.RequestedAt,
		ride.RideType, ride.PassengerCount,
		ride.EstimatedFare, ride.PaymentMethod,
		ride.PaymentStatus, ride.PickupAddress, ride.DropoffAddress,
		ride.CreatedAt, ride.UpdatedAt,
	).Scan()

	if err != nil {
		return fmt.Errorf("failed to create ride: %w", err)
	}

	return nil
}

func (r *RideRepository) Update(ctx context.Context, ride *entities.Ride) error {
	query := `
		UPDATE rides SET
			driver_id = $1, status = $2, accepted_at = $3, picked_up_at = $4, started_at = $5,
			completed_at = $6, cancelled_at = $7, cancel_reason = $8, actual_fare = $9,
			payment_status = $10, rider_rating = $11, driver_rating = $12, actual_distance = $13,
			actual_duration = $14, updated_at = $15
		WHERE id = $16
	`

	result, err := r.pool.Exec(ctx, query,
		ride.DriverID, string(ride.Status),
		ride.AcceptedAt, ride.PickedUpAt, ride.StartedAt,
		ride.CompletedAt, ride.CancelledAt, ride.CancelReason,
		ride.ActualFare, ride.PaymentStatus,
		ride.RiderRating, ride.DriverRating,
		ride.ActualDistance, ride.ActualDuration,
		ride.UpdatedAt, ride.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update ride: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ride not found: %s", ride.ID)
	}

	return nil
}

func (r *RideRepository) GetByID(ctx context.Context, id string) (*entities.Ride, error) {
	query := `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, dropoff_latitude,
		       dropoff_longitude, estimated_distance, estimated_duration, status, requested_at,
		       accepted_at, picked_up_at, started_at, completed_at, cancelled_at, cancel_reason,
		       ride_type, passenger_count, estimated_fare, actual_fare, payment_method,
		       payment_status, rider_rating, driver_rating, pickup_address, dropoff_address,
		       actual_distance, actual_duration, created_at, updated_at
		FROM rides WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanRide(ctx, query, id)
}

func (r *RideRepository) GetByRiderID(ctx context.Context, riderID string, limit int) ([]*entities.Ride, error) {
	query := `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, dropoff_latitude,
		       dropoff_longitude, estimated_distance, estimated_duration, status, requested_at,
		       accepted_at, picked_up_at, started_at, completed_at, cancelled_at, cancel_reason,
		       ride_type, passenger_count, estimated_fare, actual_fare, payment_method,
		       payment_status, rider_rating, driver_rating, pickup_address, dropoff_address,
		       actual_distance, actual_duration, created_at, updated_at
		FROM rides WHERE rider_id = $1 AND deleted_at IS NULL
		ORDER BY requested_at DESC LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, riderID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query rides: %w", err)
	}
	defer rows.Close()

	return r.scanRides(rows)
}

func (r *RideRepository) GetByDriverID(ctx context.Context, driverID string, limit int) ([]*entities.Ride, error) {
	query := `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, dropoff_latitude,
		       dropoff_longitude, estimated_distance, estimated_duration, status, requested_at,
		       accepted_at, picked_up_at, started_at, completed_at, cancelled_at, cancel_reason,
		       ride_type, passenger_count, estimated_fare, actual_fare, payment_method,
		       payment_status, rider_rating, driver_rating, pickup_address, dropoff_address,
		       actual_distance, actual_duration, created_at, updated_at
		FROM rides WHERE driver_id = $1 AND deleted_at IS NULL
		ORDER BY requested_at DESC LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, driverID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query rides: %w", err)
	}
	defer rows.Close()

	return r.scanRides(rows)
}

func (r *RideRepository) GetRequestedRides(ctx context.Context, limit int) ([]*entities.Ride, error) {
	query := `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, dropoff_latitude,
		       dropoff_longitude, estimated_distance, estimated_duration, status, requested_at,
		       accepted_at, picked_up_at, started_at, completed_at, cancelled_at, cancel_reason,
		       ride_type, passenger_count, estimated_fare, actual_fare, payment_method,
		       payment_status, rider_rating, driver_rating, pickup_address, dropoff_address,
		       actual_distance, actual_duration, created_at, updated_at
		FROM rides WHERE status = 'requested' AND deleted_at IS NULL
		ORDER BY requested_at ASC LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query requested rides: %w", err)
	}
	defer rows.Close()

	return r.scanRides(rows)
}

func (r *RideRepository) scanRide(ctx context.Context, query string, args ...interface{}) (*entities.Ride, error) {
	row := r.pool.QueryRow(ctx, query, args...)

	var (
		id, riderID, pickupAddr, dropoffAddr, rideType, paymentMethod string
		status, cancelReason, paymentStatus string
		driverID *string
		pickupLat, pickupLng, dropoffLat, dropoffLng, estimatedDist, estimatedFare, actualFare float64
		estimatedDurationMin float64
		actualDistancePtr *float64
		actualDurationMin *float64
		passengerCount int
		riderRating, driverRating *int
		requestedAt, acceptedAt, pickedUpAt, startedAt, completedAt, cancelledAt, createdAt, updatedAt time.Time
		acceptedAtPtr, pickedUpAtPtr, startedAtPtr, completedAtPtr, cancelledAtPtr *time.Time
	)

	err := row.Scan(
		&id, &riderID, &driverID,
		&pickupLat, &pickupLng, &dropoffLat, &dropoffLng,
		&estimatedDist, &estimatedDurationMin,
		&status, &requestedAt,
		&acceptedAtPtr, &pickedUpAtPtr, &startedAtPtr, &completedAtPtr, &cancelledAtPtr,
		&cancelReason, &rideType, &passengerCount, &estimatedFare, &actualFare,
		&paymentMethod, &paymentStatus, &riderRating, &driverRating,
		&pickupAddr, &dropoffAddr,
		&actualDistancePtr, &actualDurationMin,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ride not found")
		}
		return nil, fmt.Errorf("failed to scan ride: %w", err)
	}

	ride := &entities.Ride{
		ID:                   id,
		RiderID:              riderID,
		DriverID:             driverID,
		PickupLatitude:       pickupLat,
		PickupLongitude:      pickupLng,
		DropoffLatitude:      dropoffLat,
		DropoffLongitude:     dropoffLng,
		EstimatedDistance:    estimatedDist,
		EstimatedDuration:    time.Duration(estimatedDurationMin) * time.Minute,
		Status:               entities.RideStatus(status),
		RequestedAt:          requestedAt,
		AcceptedAt:           acceptedAtPtr,
		PickedUpAt:           pickedUpAtPtr,
		StartedAt:            startedAtPtr,
		CompletedAt:          completedAtPtr,
		CancelledAt:          cancelledAtPtr,
		CancelReason:         cancelReason,
		RideType:             rideType,
		PassengerCount:       passengerCount,
		EstimatedFare:        estimatedFare,
		ActualFare:           &actualFare,
		PaymentMethod:        paymentMethod,
		PaymentStatus:        paymentStatus,
		RiderRating:          riderRating,
		DriverRating:         driverRating,
		PickupAddress:        pickupAddr,
		DropoffAddress:       dropoffAddr,
		ActualDistance:       actualDistancePtr,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}

	if actualDurationMin != nil {
		duration := time.Duration(*actualDurationMin) * time.Minute
		ride.ActualDuration = &duration
	}

	return ride, nil
}

func (r *RideRepository) scanRides(rows pgx.Rows) ([]*entities.Ride, error) {
	var rides []*entities.Ride

	for rows.Next() {
		var (
			id, riderID, pickupAddr, dropoffAddr, rideType, paymentMethod string
			status, cancelReason, paymentStatus string
			driverID *string
			pickupLat, pickupLng, dropoffLat, dropoffLng, estimatedDist, estimatedFare, actualFare float64
			estimatedDurationMin float64
			actualDistancePtr *float64
			actualDurationMin *float64
			passengerCount int
			riderRating, driverRating *int
			requestedAt, acceptedAt, pickedUpAt, startedAt, completedAt, cancelledAt, createdAt, updatedAt time.Time
			acceptedAtPtr, pickedUpAtPtr, startedAtPtr, completedAtPtr, cancelledAtPtr *time.Time
		)

		err := rows.Scan(
			&id, &riderID, &driverID,
			&pickupLat, &pickupLng, &dropoffLat, &dropoffLng,
			&estimatedDist, &estimatedDurationMin,
			&status, &requestedAt,
			&acceptedAtPtr, &pickedUpAtPtr, &startedAtPtr, &completedAtPtr, &cancelledAtPtr,
			&cancelReason, &rideType, &passengerCount, &estimatedFare, &actualFare,
			&paymentMethod, &paymentStatus, &riderRating, &driverRating,
			&pickupAddr, &dropoffAddr,
			&actualDistancePtr, &actualDurationMin,
			&createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan ride: %w", err)
		}

		ride := &entities.Ride{
			ID:                   id,
			RiderID:              riderID,
			DriverID:             driverID,
			PickupLatitude:       pickupLat,
			PickupLongitude:      pickupLng,
			DropoffLatitude:      dropoffLat,
			DropoffLongitude:     dropoffLng,
			EstimatedDistance:    estimatedDist,
			EstimatedDuration:    time.Duration(estimatedDurationMin) * time.Minute,
			Status:               entities.RideStatus(status),
			RequestedAt:          requestedAt,
			AcceptedAt:           acceptedAtPtr,
			PickedUpAt:           pickedUpAtPtr,
			StartedAt:            startedAtPtr,
			CompletedAt:          completedAtPtr,
			CancelledAt:          cancelledAtPtr,
			CancelReason:         cancelReason,
			RideType:             rideType,
			PassengerCount:       passengerCount,
			EstimatedFare:        estimatedFare,
			ActualFare:           &actualFare,
			PaymentMethod:        paymentMethod,
			PaymentStatus:        paymentStatus,
			RiderRating:          riderRating,
			DriverRating:         driverRating,
			PickupAddress:        pickupAddr,
			DropoffAddress:       dropoffAddr,
			ActualDistance:       actualDistancePtr,
			CreatedAt:            createdAt,
			UpdatedAt:            updatedAt,
		}

		if actualDurationMin != nil {
			duration := time.Duration(*actualDurationMin) * time.Minute
			ride.ActualDuration = &duration
		}

		rides = append(rides, ride)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rides: %w", err)
	}

	return rides, nil
}
