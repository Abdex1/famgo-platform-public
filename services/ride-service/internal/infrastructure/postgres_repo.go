// services/ride-service/internal/infrastructure/postgres_repo.go
// PostgreSQL Ride Repository Implementations

package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// PostgresRideRepository implements RideRepository
type PostgresRideRepository struct {
	db *sql.DB
}

func NewPostgresRideRepository(db *sql.DB) *PostgresRideRepository {
	return &PostgresRideRepository{db: db}
}

// ===== RIDE REPOSITORY METHODS =====

func (r *PostgresRideRepository) GetRide(ctx context.Context, rideID string) (*domain.Ride, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time, 
                cancellation_reason, created_at, updated_at
         FROM rides WHERE id = $1`,
		rideID)

	ride := &domain.Ride{}
	err := row.Scan(
		&ride.ID,
		&ride.PassengerID,
		&ride.DriverID,
		&ride.PickupLat,
		&ride.PickupLon,
		&ride.DropoffLat,
		&ride.DropoffLon,
		&ride.Status,
		&ride.EstimatedFare,
		&ride.ActualFare,
		&ride.PickupTime,
		&ride.DropoffTime,
		&ride.CancellationReason,
		&ride.CreatedAt,
		&ride.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRideNotFound
		}
		return nil, err
	}

	return ride, nil
}

func (r *PostgresRideRepository) CreateRide(ctx context.Context, ride *domain.Ride) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO rides (id, passenger_id, driver_id, pickup_lat, pickup_lon, 
                          dropoff_lat, dropoff_lon, status, estimated_fare, actual_fare, 
                          pickup_time, dropoff_time, cancellation_reason, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		ride.ID,
		ride.PassengerID,
		ride.DriverID,
		ride.PickupLat,
		ride.PickupLon,
		ride.DropoffLat,
		ride.DropoffLon,
		ride.Status,
		ride.EstimatedFare,
		ride.ActualFare,
		ride.PickupTime,
		ride.DropoffTime,
		ride.CancellationReason,
		ride.CreatedAt,
		ride.UpdatedAt,
	)
	return err
}

func (r *PostgresRideRepository) UpdateRide(ctx context.Context, ride *domain.Ride) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE rides SET passenger_id = $1, driver_id = $2, pickup_lat = $3, pickup_lon = $4,
                         dropoff_lat = $5, dropoff_lon = $6, status = $7, estimated_fare = $8,
                         actual_fare = $9, pickup_time = $10, dropoff_time = $11,
                         cancellation_reason = $12, updated_at = $13
         WHERE id = $14`,
		ride.PassengerID,
		ride.DriverID,
		ride.PickupLat,
		ride.PickupLon,
		ride.DropoffLat,
		ride.DropoffLon,
		ride.Status,
		ride.EstimatedFare,
		ride.ActualFare,
		ride.PickupTime,
		ride.DropoffTime,
		ride.CancellationReason,
		ride.UpdatedAt,
		ride.ID,
	)
	return err
}

func (r *PostgresRideRepository) GetRidesByPassenger(ctx context.Context, passengerID string, limit, offset int) ([]domain.Ride, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time,
                cancellation_reason, created_at, updated_at
         FROM rides WHERE passenger_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		passengerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []domain.Ride
	for rows.Next() {
		ride := domain.Ride{}
		err := rows.Scan(
			&ride.ID,
			&ride.PassengerID,
			&ride.DriverID,
			&ride.PickupLat,
			&ride.PickupLon,
			&ride.DropoffLat,
			&ride.DropoffLon,
			&ride.Status,
			&ride.EstimatedFare,
			&ride.ActualFare,
			&ride.PickupTime,
			&ride.DropoffTime,
			&ride.CancellationReason,
			&ride.CreatedAt,
			&ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, rows.Err()
}

func (r *PostgresRideRepository) GetRidesByDriver(ctx context.Context, driverID string, limit, offset int) ([]domain.Ride, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time,
                cancellation_reason, created_at, updated_at
         FROM rides WHERE driver_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		driverID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []domain.Ride
	for rows.Next() {
		ride := domain.Ride{}
		err := rows.Scan(
			&ride.ID,
			&ride.PassengerID,
			&ride.DriverID,
			&ride.PickupLat,
			&ride.PickupLon,
			&ride.DropoffLat,
			&ride.DropoffLon,
			&ride.Status,
			&ride.EstimatedFare,
			&ride.ActualFare,
			&ride.PickupTime,
			&ride.DropoffTime,
			&ride.CancellationReason,
			&ride.CreatedAt,
			&ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, rows.Err()
}

func (r *PostgresRideRepository) GetActiveRides(ctx context.Context) ([]domain.Ride, error) {
	activeStatuses := []string{
		string(domain.RideStatusRequested),
		string(domain.RideStatusSearching),
		string(domain.RideStatusAssigned),
		string(domain.RideStatusDriverArriving),
		string(domain.RideStatusStarted),
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time,
                cancellation_reason, created_at, updated_at
         FROM rides WHERE status = ANY($1) ORDER BY created_at DESC`,
		activeStatuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []domain.Ride
	for rows.Next() {
		ride := domain.Ride{}
		err := rows.Scan(
			&ride.ID,
			&ride.PassengerID,
			&ride.DriverID,
			&ride.PickupLat,
			&ride.PickupLon,
			&ride.DropoffLat,
			&ride.DropoffLon,
			&ride.Status,
			&ride.EstimatedFare,
			&ride.ActualFare,
			&ride.PickupTime,
			&ride.DropoffTime,
			&ride.CancellationReason,
			&ride.CreatedAt,
			&ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, rows.Err()
}

// ===== RIDE STATUS HISTORY REPOSITORY =====

// PostgresRideStatusHistoryRepository for status transitions
type PostgresRideStatusHistoryRepository struct {
	db *sql.DB
}

func NewPostgresRideStatusHistoryRepository(db *sql.DB) *PostgresRideStatusHistoryRepository {
	return &PostgresRideStatusHistoryRepository{db: db}
}

func (r *PostgresRideStatusHistoryRepository) CreateStatusHistory(ctx context.Context, history *domain.RideStatusHistory) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ride_status_history (id, ride_id, old_status, new_status, changed_at)
         VALUES ($1, $2, $3, $4, $5)`,
		history.ID,
		history.RideID,
		history.OldStatus,
		history.NewStatus,
		history.ChangedAt,
	)
	return err
}

func (r *PostgresRideStatusHistoryRepository) GetRideStatusHistory(ctx context.Context, rideID string) ([]domain.RideStatusHistory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ride_id, old_status, new_status, changed_at
         FROM ride_status_history WHERE ride_id = $1 ORDER BY changed_at ASC`,
		rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []domain.RideStatusHistory
	for rows.Next() {
		history := domain.RideStatusHistory{}
		err := rows.Scan(
			&history.ID,
			&history.RideID,
			&history.OldStatus,
			&history.NewStatus,
			&history.ChangedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, rows.Err()
}

func (r *PostgresRideStatusHistoryRepository) GetLatestStatusTransition(ctx context.Context, rideID string) (*domain.RideStatusHistory, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, ride_id, old_status, new_status, changed_at
         FROM ride_status_history WHERE ride_id = $1 ORDER BY changed_at DESC LIMIT 1`,
		rideID)

	history := &domain.RideStatusHistory{}
	err := row.Scan(
		&history.ID,
		&history.RideID,
		&history.OldStatus,
		&history.NewStatus,
		&history.ChangedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return history, nil
}

// GetRidesByStatusAfterTime gets rides with specific status after a timestamp (for monitoring)
func (r *PostgresRideRepository) GetRidesByStatusAfterTime(ctx context.Context, status domain.RideStatus, afterTime time.Time) ([]domain.Ride, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time,
                cancellation_reason, created_at, updated_at
         FROM rides WHERE status = $1 AND created_at > $2 ORDER BY created_at DESC`,
		status, afterTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []domain.Ride
	for rows.Next() {
		ride := domain.Ride{}
		err := rows.Scan(
			&ride.ID,
			&ride.PassengerID,
			&ride.DriverID,
			&ride.PickupLat,
			&ride.PickupLon,
			&ride.DropoffLat,
			&ride.DropoffLon,
			&ride.Status,
			&ride.EstimatedFare,
			&ride.ActualFare,
			&ride.PickupTime,
			&ride.DropoffTime,
			&ride.CancellationReason,
			&ride.CreatedAt,
			&ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, rows.Err()
}

// GetRidesByStatusBetweenTime gets rides for analytics queries
func (r *PostgresRideRepository) GetRidesByStatusBetweenTime(ctx context.Context, status domain.RideStatus, startTime, endTime time.Time) ([]domain.Ride, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, passenger_id, driver_id, pickup_lat, pickup_lon, dropoff_lat, dropoff_lon,
                status, estimated_fare, actual_fare, pickup_time, dropoff_time,
                cancellation_reason, created_at, updated_at
         FROM rides WHERE status = $1 AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC`,
		status, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []domain.Ride
	for rows.Next() {
		ride := domain.Ride{}
		err := rows.Scan(
			&ride.ID,
			&ride.PassengerID,
			&ride.DriverID,
			&ride.PickupLat,
			&ride.PickupLon,
			&ride.DropoffLat,
			&ride.DropoffLon,
			&ride.Status,
			&ride.EstimatedFare,
			&ride.ActualFare,
			&ride.PickupTime,
			&ride.DropoffTime,
			&ride.CancellationReason,
			&ride.CreatedAt,
			&ride.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, rows.Err()
}
