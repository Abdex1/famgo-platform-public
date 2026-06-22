package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/entities"
	"github.com/lib/pq"
)

// RideRepository handles ride data persistence
type RideRepository struct {
	db *sql.DB
}

// NewRideRepository creates new ride repository
func NewRideRepository(db *sql.DB) *RideRepository {
	return &RideRepository{db: db}
}

// CreateRideRequest saves new ride request
func (r *RideRepository) CreateRideRequest(ctx context.Context, req *entities.RideRequest) error {
	query := `
		INSERT INTO ride_requests 
		(id, rider_id, pickup_lat, pickup_lng, pickup_address, dropoff_lat, dropoff_lng, dropoff_address, 
		 ride_type, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.RiderID, req.PickupLat, req.PickupLng, req.PickupAddress,
		req.DropoffLat, req.DropoffLng, req.DropoffAddress,
		req.RideType, req.Status, req.CreatedAt, req.ExpiresAt,
	)
	return err
}

// GetRideRequest retrieves ride request by ID
func (r *RideRepository) GetRideRequest(ctx context.Context, requestID string) (*entities.RideRequest, error) {
	query := `
		SELECT id, rider_id, pickup_lat, pickup_lng, pickup_address, dropoff_lat, dropoff_lng, 
		       dropoff_address, ride_type, status, created_at, expires_at
		FROM ride_requests WHERE id = $1
	`
	
	req := &entities.RideRequest{}
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.ID, &req.RiderID, &req.PickupLat, &req.PickupLng, &req.PickupAddress,
		&req.DropoffLat, &req.DropoffLng, &req.DropoffAddress,
		&req.RideType, &req.Status, &req.CreatedAt, &req.ExpiresAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("ride request not found")
	}
	return req, err
}

// CreateRide saves assigned ride
func (r *RideRepository) CreateRide(ctx context.Context, ride *entities.Ride) error {
	query := `
		INSERT INTO rides 
		(id, request_id, rider_id, driver_id, pickup_lat, pickup_lng, pickup_address, 
		 dropoff_lat, dropoff_lng, dropoff_address, ride_type, status, surge_multiplier, 
		 estimated_fare, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		ride.ID, ride.RequestID, ride.RiderID, ride.DriverID, ride.PickupLat, ride.PickupLng,
		ride.PickupAddress, ride.DropoffLat, ride.DropoffLng, ride.DropoffAddress,
		ride.RideType, ride.Status, ride.SurgeMultiplier, ride.EstimatedFare,
		ride.CreatedAt, ride.UpdatedAt,
	)
	return err
}

// GetRide retrieves ride by ID
func (r *RideRepository) GetRide(ctx context.Context, rideID string) (*entities.Ride, error) {
	query := `
		SELECT id, request_id, rider_id, driver_id, pickup_lat, pickup_lng, pickup_address,
		       dropoff_lat, dropoff_lng, dropoff_address, ride_type, status, actual_distance_meters,
		       actual_duration_seconds, actual_fare, surge_multiplier, estimated_fare, pickup_time,
		       dropoff_time, cancellation_reason, cancelled_by, assigned_at, accepted_at, 
		       arrived_at, started_at, completed_at, created_at, updated_at, payment_method_id,
		       promotion_code_id
		FROM rides WHERE id = $1
	`
	
	ride := &entities.Ride{}
	err := r.db.QueryRowContext(ctx, query, rideID).Scan(
		&ride.ID, &ride.RequestID, &ride.RiderID, &ride.DriverID, &ride.PickupLat, &ride.PickupLng,
		&ride.PickupAddress, &ride.DropoffLat, &ride.DropoffLng, &ride.DropoffAddress,
		&ride.RideType, &ride.Status, &ride.ActualDistance, &ride.ActualDuration, &ride.ActualFare,
		&ride.SurgeMultiplier, &ride.EstimatedFare, &ride.PickupTime, &ride.DropoffTime,
		&ride.CancellationReason, &ride.CancelledBy, &ride.AssignedAt, &ride.AcceptedAt,
		&ride.ArrivedAt, &ride.StartedAt, &ride.CompletedAt, &ride.CreatedAt, &ride.UpdatedAt,
		&ride.PaymentMethodID, &ride.PromotionCodeID,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("ride not found")
	}
	return ride, err
}

// UpdateRideStatus updates ride status
func (r *RideRepository) UpdateRideStatus(ctx context.Context, rideID string, status entities.RideStatus) error {
	query := `UPDATE rides SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), rideID)
	return err
}

// UpdateRideWithDriver updates ride with assigned driver
func (r *RideRepository) UpdateRideWithDriver(ctx context.Context, rideID, driverID string) error {
	query := `
		UPDATE rides 
		SET driver_id = $1, status = $2, assigned_at = $3, updated_at = $4 
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, driverID, entities.StatusMatched, time.Now(), time.Now(), rideID)
	return err
}

// SaveRideLocation stores GPS location point
func (r *RideRepository) SaveRideLocation(ctx context.Context, loc *entities.RideLocation) error {
	query := `
		INSERT INTO ride_locations (id, ride_id, latitude, longitude, heading, speed, accuracy, source, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		loc.ID, loc.RideID, loc.Latitude, loc.Longitude, loc.Heading, loc.Speed, loc.Accuracy, loc.Source, loc.RecordedAt,
	)
	return err
}

// GetRideLocations retrieves all GPS points for ride
func (r *RideRepository) GetRideLocations(ctx context.Context, rideID string) ([]entities.RideLocation, error) {
	query := `
		SELECT id, ride_id, latitude, longitude, heading, speed, accuracy, source, recorded_at
		FROM ride_locations WHERE ride_id = $1 ORDER BY recorded_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var locations []entities.RideLocation
	for rows.Next() {
		var loc entities.RideLocation
		err := rows.Scan(&loc.ID, &loc.RideID, &loc.Latitude, &loc.Longitude, &loc.Heading, &loc.Speed, &loc.Accuracy, &loc.Source, &loc.RecordedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	
	return locations, rows.Err()
}

// GetRiderRideHistory retrieves ride history for rider
func (r *RideRepository) GetRiderRideHistory(ctx context.Context, riderID string, limit, offset int) ([]entities.Ride, error) {
	query := `
		SELECT id, request_id, rider_id, driver_id, pickup_lat, pickup_lng, pickup_address,
		       dropoff_lat, dropoff_lng, dropoff_address, ride_type, status, actual_distance_meters,
		       actual_duration_seconds, actual_fare, surge_multiplier, estimated_fare, created_at, updated_at
		FROM rides 
		WHERE rider_id = $1 AND status IN ('COMPLETED', 'CANCELLED', 'NO_SHOW')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, riderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rides []entities.Ride
	for rows.Next() {
		var ride entities.Ride
		err := rows.Scan(
			&ride.ID, &ride.RequestID, &ride.RiderID, &ride.DriverID, &ride.PickupLat, &ride.PickupLng,
			&ride.PickupAddress, &ride.DropoffLat, &ride.DropoffLng, &ride.DropoffAddress,
			&ride.RideType, &ride.Status, &ride.ActualDistance, &ride.ActualDuration, &ride.ActualFare,
			&ride.SurgeMultiplier, &ride.EstimatedFare, &ride.CreatedAt, &ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	
	return rides, rows.Err()
}

// CreateRideSession creates driver-rider session
func (r *RideRepository) CreateRideSession(ctx context.Context, session *entities.RideSession) error {
	query := `
		INSERT INTO ride_sessions (id, ride_id, driver_id, rider_id, status, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.RideID, session.DriverID, session.RiderID, session.Status, session.StartedAt, session.CreatedAt, session.UpdatedAt,
	)
	return err
}

// GetActiveRidesByDriver retrieves active rides for driver
func (r *RideRepository) GetActiveRidesByDriver(ctx context.Context, driverID string) ([]entities.Ride, error) {
	query := `
		SELECT id, request_id, rider_id, driver_id, pickup_lat, pickup_lng, pickup_address,
		       dropoff_lat, dropoff_lng, dropoff_address, ride_type, status, created_at, updated_at
		FROM rides 
		WHERE driver_id = $1 AND status IN ('ACCEPTED', 'ARRIVED', 'IN_PROGRESS')
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rides []entities.Ride
	for rows.Next() {
		var ride entities.Ride
		err := rows.Scan(
			&ride.ID, &ride.RequestID, &ride.RiderID, &ride.DriverID, &ride.PickupLat, &ride.PickupLng,
			&ride.PickupAddress, &ride.DropoffLat, &ride.DropoffLng, &ride.DropoffAddress,
			&ride.RideType, &ride.Status, &ride.CreatedAt, &ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	
	return rides, rows.Err()
}
