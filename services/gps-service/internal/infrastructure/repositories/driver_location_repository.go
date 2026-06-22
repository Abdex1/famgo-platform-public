// services/gps-service/internal/infrastructure/repositories/driver_location_repository.go
// PostgreSQL repository for driver location persistence

package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// DriverLocationRepository handles driver location data in PostgreSQL
type DriverLocationRepository struct {
	pool *pgxpool.Pool
}

// NewDriverLocationRepository creates a new driver location repository
func NewDriverLocationRepository(pool *pgxpool.Pool) *DriverLocationRepository {
	return &DriverLocationRepository{
		pool: pool,
	}
}

// Create inserts a new driver location record
func (r *DriverLocationRepository) Create(ctx context.Context, dl *entities.DriverLocation) error {
	query := `
		INSERT INTO driver_locations (
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17,
			$18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30
		)
	`

	err := r.pool.QueryRow(ctx, query,
		dl.ID, dl.DriverID,
		dl.CurrentLocation.Coordinates.Latitude,
		dl.CurrentLocation.Coordinates.Longitude,
		dl.CurrentLocation.Altitude,
		dl.CurrentLocation.Accuracy,
		dl.CurrentLocation.Speed,
		dl.CurrentLocation.Heading,
		string(dl.Status), dl.IsOnline, dl.LastUpdateAt, dl.LastSyncAt, dl.LastSeenAt,
		dl.ConsecutiveFailures, dl.VehicleID, dl.VehicleRegistration, dl.ServiceStatus,
		dl.AcceptedRideCount, dl.CompletedRideCount, dl.CancelledRideCount,
		dl.AverageRating, dl.RecentAcceptanceRate, dl.GeohashPrefix, dl.LastLocationHash,
		dl.IsVerified, dl.IsDocumentsExpired, dl.IsBanned, dl.CreatedAt, dl.UpdatedAt, dl.CreatedBy,
	).Scan()

	if err != nil {
		return fmt.Errorf("failed to create driver location: %w", err)
	}

	return nil
}

// Update updates an existing driver location record
func (r *DriverLocationRepository) Update(ctx context.Context, dl *entities.DriverLocation) error {
	query := `
		UPDATE driver_locations SET
			latitude = $1, longitude = $2, altitude = $3, accuracy = $4,
			speed = $5, heading = $6, status = $7, is_online = $8,
			last_update_at = $9, last_sync_at = $10, last_seen_at = $11,
			consecutive_failures = $12, vehicle_id = $13, service_status = $14,
			accepted_ride_count = $15, completed_ride_count = $16, cancelled_ride_count = $17,
			average_rating = $18, recent_acceptance_rate = $19, geohash_prefix = $20,
			last_location_hash = $21, is_verified = $22, is_documents_expired = $23,
			is_banned = $24, updated_at = $25, updated_by = $26
		WHERE driver_id = $27
	`

	result, err := r.pool.Exec(ctx, query,
		dl.CurrentLocation.Coordinates.Latitude,
		dl.CurrentLocation.Coordinates.Longitude,
		dl.CurrentLocation.Altitude,
		dl.CurrentLocation.Accuracy,
		dl.CurrentLocation.Speed,
		dl.CurrentLocation.Heading,
		string(dl.Status), dl.IsOnline,
		dl.LastUpdateAt, dl.LastSyncAt, dl.LastSeenAt,
		dl.ConsecutiveFailures, dl.VehicleID, dl.ServiceStatus,
		dl.AcceptedRideCount, dl.CompletedRideCount, dl.CancelledRideCount,
		dl.AverageRating, dl.RecentAcceptanceRate, dl.GeohashPrefix,
		dl.LastLocationHash, dl.IsVerified, dl.IsDocumentsExpired,
		dl.IsBanned, dl.UpdatedAt, dl.UpdatedBy,
		dl.DriverID,
	)

	if err != nil {
		return fmt.Errorf("failed to update driver location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("driver location not found: %s", dl.DriverID)
	}

	return nil
}

// GetByDriverID retrieves a driver location by driver ID
func (r *DriverLocationRepository) GetByDriverID(ctx context.Context, driverID string) (*entities.DriverLocation, error) {
	query := `
		SELECT
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by, updated_by
		FROM driver_locations
		WHERE driver_id = $1 AND deleted_at IS NULL
	`

	return r.scanDriverLocation(ctx, query, driverID)
}

// GetByID retrieves a driver location by ID
func (r *DriverLocationRepository) GetByID(ctx context.Context, id string) (*entities.DriverLocation, error) {
	query := `
		SELECT
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by, updated_by
		FROM driver_locations
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanDriverLocation(ctx, query, id)
}

// GetOnlineDrivers retrieves all online drivers
func (r *DriverLocationRepository) GetOnlineDrivers(ctx context.Context) ([]*entities.DriverLocation, error) {
	query := `
		SELECT
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by, updated_by
		FROM driver_locations
		WHERE is_online = true AND status = 'online' AND deleted_at IS NULL
		ORDER BY last_update_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query online drivers: %w", err)
	}
	defer rows.Close()

	return r.scanDriverLocations(rows)
}

// GetNearbyDrivers retrieves drivers within bounds
func (r *DriverLocationRepository) GetNearbyDrivers(
	ctx context.Context,
	minLat, maxLat, minLng, maxLng float64,
	limit int,
) ([]*entities.DriverLocation, error) {
	query := `
		SELECT
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by, updated_by
		FROM driver_locations
		WHERE latitude >= $1 AND latitude <= $2
		  AND longitude >= $3 AND longitude <= $4
		  AND is_online = true AND deleted_at IS NULL
		ORDER BY last_update_at DESC
		LIMIT $5
	`

	rows, err := r.pool.Query(ctx, query, minLat, maxLat, minLng, maxLng, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby drivers: %w", err)
	}
	defer rows.Close()

	return r.scanDriverLocations(rows)
}

// GetStaleLocations retrieves locations older than specified duration
func (r *DriverLocationRepository) GetStaleLocations(
	ctx context.Context,
	staleDuration time.Duration,
) ([]*entities.DriverLocation, error) {
	query := `
		SELECT
			id, driver_id, latitude, longitude, altitude, accuracy, speed, heading,
			status, is_online, last_update_at, last_sync_at, last_seen_at,
			consecutive_failures, vehicle_id, vehicle_registration, service_status,
			accepted_ride_count, completed_ride_count, cancelled_ride_count,
			average_rating, recent_acceptance_rate, geohash_prefix, last_location_hash,
			is_verified, is_documents_expired, is_banned, created_at, updated_at, created_by, updated_by
		FROM driver_locations
		WHERE last_update_at < NOW() - INTERVAL '1 second' * $1
		  AND deleted_at IS NULL
		ORDER BY last_update_at ASC
	`

	rows, err := r.pool.Query(ctx, query, staleDuration.Seconds())
	if err != nil {
		return nil, fmt.Errorf("failed to query stale locations: %w", err)
	}
	defer rows.Close()

	return r.scanDriverLocations(rows)
}

// Delete marks a driver location as deleted
func (r *DriverLocationRepository) Delete(ctx context.Context, driverID string) error {
	query := `UPDATE driver_locations SET deleted_at = NOW() WHERE driver_id = $1`
	result, err := r.pool.Exec(ctx, query, driverID)
	if err != nil {
		return fmt.Errorf("failed to delete driver location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("driver location not found: %s", driverID)
	}

	return nil
}

// Helper functions

func (r *DriverLocationRepository) scanDriverLocation(ctx context.Context, query string, args ...interface{}) (*entities.DriverLocation, error) {
	row := r.pool.QueryRow(ctx, query, args...)

	var (
		id, driverID, vehicleID, vehicleReg, serviceStatus string
		status, geohashPrefix, lastLocationHash            string
		latitude, longitude, altitude, accuracy, speed, heading float64
		acceptedRides, completedRides, cancelledRides, consecutiveFailures int
		averageRating, acceptanceRate float64
		isOnline, isVerified, isDocumentsExpired, isBanned bool
		lastUpdateAt, lastSyncAt, lastSeenAt, createdAt, updatedAt time.Time
		createdBy, updatedBy string
	)

	err := row.Scan(
		&id, &driverID, &latitude, &longitude, &altitude, &accuracy, &speed, &heading,
		&status, &isOnline, &lastUpdateAt, &lastSyncAt, &lastSeenAt,
		&consecutiveFailures, &vehicleID, &vehicleReg, &serviceStatus,
		&acceptedRides, &completedRides, &cancelledRides,
		&averageRating, &acceptanceRate, &geohashPrefix, &lastLocationHash,
		&isVerified, &isDocumentsExpired, &isBanned, &createdAt, &updatedAt, &createdBy, &updatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("driver location not found")
		}
		return nil, fmt.Errorf("failed to scan driver location: %w", err)
	}

	location, _ := valueobjects.NewGeolocation(latitude, longitude, altitude, accuracy, speed, heading, 0)

	dl := &entities.DriverLocation{
		ID:                  id,
		DriverID:            driverID,
		CurrentLocation:     location,
		Status:              entities.DriverStatus(status),
		IsOnline:            isOnline,
		LastUpdateAt:        lastUpdateAt,
		LastSyncAt:          lastSyncAt,
		LastSeenAt:          lastSeenAt,
		ConsecutiveFailures: consecutiveFailures,
		VehicleID:           vehicleID,
		VehicleRegistration: vehicleReg,
		ServiceStatus:       serviceStatus,
		AcceptedRideCount:   acceptedRides,
		CompletedRideCount:  completedRides,
		CancelledRideCount:  cancelledRides,
		AverageRating:       averageRating,
		RecentAcceptanceRate: acceptanceRate,
		GeohashPrefix:       geohashPrefix,
		LastLocationHash:    lastLocationHash,
		IsVerified:          isVerified,
		IsDocumentsExpired:  isDocumentsExpired,
		IsBanned:            isBanned,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		CreatedBy:           createdBy,
		UpdatedBy:           updatedBy,
	}

	return dl, nil
}

func (r *DriverLocationRepository) scanDriverLocations(rows pgx.Rows) ([]*entities.DriverLocation, error) {
	var locations []*entities.DriverLocation

	for rows.Next() {
		var (
			id, driverID, vehicleID, vehicleReg, serviceStatus string
			status, geohashPrefix, lastLocationHash            string
			latitude, longitude, altitude, accuracy, speed, heading float64
			acceptedRides, completedRides, cancelledRides, consecutiveFailures int
			averageRating, acceptanceRate float64
			isOnline, isVerified, isDocumentsExpired, isBanned bool
			lastUpdateAt, lastSyncAt, lastSeenAt, createdAt, updatedAt time.Time
			createdBy, updatedBy string
		)

		err := rows.Scan(
			&id, &driverID, &latitude, &longitude, &altitude, &accuracy, &speed, &heading,
			&status, &isOnline, &lastUpdateAt, &lastSyncAt, &lastSeenAt,
			&consecutiveFailures, &vehicleID, &vehicleReg, &serviceStatus,
			&acceptedRides, &completedRides, &cancelledRides,
			&averageRating, &acceptanceRate, &geohashPrefix, &lastLocationHash,
			&isVerified, &isDocumentsExpired, &isBanned, &createdAt, &updatedAt, &createdBy, &updatedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan driver location: %w", err)
		}

		location, _ := valueobjects.NewGeolocation(latitude, longitude, altitude, accuracy, speed, heading, 0)

		dl := &entities.DriverLocation{
			ID:                   id,
			DriverID:             driverID,
			CurrentLocation:      location,
			Status:               entities.DriverStatus(status),
			IsOnline:             isOnline,
			LastUpdateAt:         lastUpdateAt,
			LastSyncAt:           lastSyncAt,
			LastSeenAt:           lastSeenAt,
			ConsecutiveFailures:  consecutiveFailures,
			VehicleID:            vehicleID,
			VehicleRegistration:  vehicleReg,
			ServiceStatus:        serviceStatus,
			AcceptedRideCount:    acceptedRides,
			CompletedRideCount:   completedRides,
			CancelledRideCount:   cancelledRides,
			AverageRating:        averageRating,
			RecentAcceptanceRate: acceptanceRate,
			GeohashPrefix:        geohashPrefix,
			LastLocationHash:     lastLocationHash,
			IsVerified:           isVerified,
			IsDocumentsExpired:   isDocumentsExpired,
			IsBanned:             isBanned,
			CreatedAt:            createdAt,
			UpdatedAt:            updatedAt,
			CreatedBy:            createdBy,
			UpdatedBy:            updatedBy,
		}

		locations = append(locations, dl)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate driver locations: %w", err)
	}

	return locations, nil
}
