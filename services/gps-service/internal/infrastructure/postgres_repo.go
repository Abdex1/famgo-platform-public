package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abdex1/FamGo-platform/packages/redis-platform"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
)

// PostgresLocationRepository implements LocationRepository using PostgreSQL
type PostgresLocationRepository struct {
	db *sql.DB
}

// NewPostgresLocationRepository creates a new PostgreSQL location repository
func NewPostgresLocationRepository(db *sql.DB) *PostgresLocationRepository {
	return &PostgresLocationRepository{db: db}
}

// GetDriverLocation retrieves current location for a driver
func (r *PostgresLocationRepository) GetDriverLocation(ctx context.Context, driverID string) (*domain.DriverLocation, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, driver_id, latitude, longitude, accuracy, updated_at
		 FROM driver_locations 
		 WHERE driver_id = $1
		 ORDER BY updated_at DESC
		 LIMIT 1`,
		driverID)

	loc := &domain.DriverLocation{}
	err := row.Scan(&loc.ID, &loc.DriverID, &loc.Latitude, &loc.Longitude, &loc.Accuracy, &loc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query location: %w", err)
	}

	return loc, nil
}

// UpdateDriverLocation updates or creates driver location
func (r *PostgresLocationRepository) UpdateDriverLocation(ctx context.Context, location *domain.DriverLocation) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO driver_locations (id, driver_id, latitude, longitude, accuracy, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (driver_id) DO UPDATE SET
			id = EXCLUDED.id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			accuracy = EXCLUDED.accuracy,
			updated_at = EXCLUDED.updated_at`,
		location.ID,
		location.DriverID,
		location.Latitude,
		location.Longitude,
		location.Accuracy,
		location.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	return nil
}

// ListActiveLocations retrieves all active driver locations
func (r *PostgresLocationRepository) ListActiveLocations(ctx context.Context) ([]domain.DriverLocation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, driver_id, latitude, longitude, accuracy, updated_at
		 FROM driver_locations
		 WHERE updated_at > NOW() - INTERVAL '5 minutes'
		 ORDER BY updated_at DESC`)

	if err != nil {
		return nil, fmt.Errorf("failed to query active locations: %w", err)
	}
	defer rows.Close()

	var locations []domain.DriverLocation
	for rows.Next() {
		loc := domain.DriverLocation{}
		err := rows.Scan(&loc.ID, &loc.DriverID, &loc.Latitude, &loc.Longitude, &loc.Accuracy, &loc.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating locations: %w", err)
	}

	return locations, nil
}

// DeleteDriverLocation removes driver location
func (r *PostgresLocationRepository) DeleteDriverLocation(ctx context.Context, driverID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM driver_locations WHERE driver_id = $1`,
		driverID)

	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	return nil
}

// ========================================

// PostgresTripRepository implements TripRepository using PostgreSQL
type PostgresTripRepository struct {
	db *sql.DB
}

// NewPostgresTripRepository creates a new PostgreSQL trip repository
func NewPostgresTripRepository(db *sql.DB) *PostgresTripRepository {
	return &PostgresTripRepository{db: db}
}

// GetTrip retrieves a trip by ID
func (r *PostgresTripRepository) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, ride_id, driver_id, started_at, status
		 FROM trips WHERE id = $1`,
		tripID)

	trip := &domain.Trip{Route: []domain.RoutePoint{}}
	err := row.Scan(&trip.ID, &trip.RideID, &trip.DriverID, &trip.StartedAt, &trip.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query trip: %w", err)
	}

	return trip, nil
}

// CreateTrip creates a new trip
func (r *PostgresTripRepository) CreateTrip(ctx context.Context, trip *domain.Trip) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO trips (id, ride_id, driver_id, started_at, status)
		 VALUES ($1, $2, $3, $4, $5)`,
		trip.ID, trip.RideID, trip.DriverID, trip.StartedAt, trip.Status)

	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

// UpdateTrip updates an existing trip
func (r *PostgresTripRepository) UpdateTrip(ctx context.Context, trip *domain.Trip) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE trips SET status = $1 WHERE id = $2`,
		trip.Status, trip.ID)

	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	return nil
}

// AddRoutePoint adds a location point to a trip's route
func (r *PostgresTripRepository) AddRoutePoint(ctx context.Context, tripID string, point domain.RoutePoint) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO trip_route_points (id, trip_id, latitude, longitude, timestamp)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4)`,
		tripID, point.Latitude, point.Longitude, point.Timestamp)

	if err != nil {
		return fmt.Errorf("failed to add route point: %w", err)
	}

	return nil
}

// GetTripsByDriver retrieves all active trips for a driver
func (r *PostgresTripRepository) GetTripsByDriver(ctx context.Context, driverID string) ([]domain.Trip, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ride_id, driver_id, started_at, status
		 FROM trips WHERE driver_id = $1 AND status = $2`,
		driverID, domain.TripStatusActive)

	if err != nil {
		return nil, fmt.Errorf("failed to query trips: %w", err)
	}
	defer rows.Close()

	var trips []domain.Trip
	for rows.Next() {
		trip := domain.Trip{Route: []domain.RoutePoint{}}
		err := rows.Scan(&trip.ID, &trip.RideID, &trip.DriverID, &trip.StartedAt, &trip.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trip: %w", err)
		}
		trips = append(trips, trip)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trips: %w", err)
	}

	return trips, nil
}

// ========================================

// PostgresGeofenceRepository implements GeofenceRepository using PostgreSQL
type PostgresGeofenceRepository struct {
	db *sql.DB
}

// NewPostgresGeofenceRepository creates a new PostgreSQL geofence repository
func NewPostgresGeofenceRepository(db *sql.DB) *PostgresGeofenceRepository {
	return &PostgresGeofenceRepository{db: db}
}

// GetGeofence retrieves a geofence by ID
func (r *PostgresGeofenceRepository) GetGeofence(ctx context.Context, geofenceID string) (*domain.Geofence, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, latitude, longitude, radius, created_at
		 FROM geofences WHERE id = $1`,
		geofenceID)

	gf := &domain.Geofence{}
	err := row.Scan(&gf.ID, &gf.Name, &gf.Latitude, &gf.Longitude, &gf.Radius, &gf.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query geofence: %w", err)
	}

	return gf, nil
}

// GetAllGeofences retrieves all geofences
func (r *PostgresGeofenceRepository) GetAllGeofences(ctx context.Context) ([]domain.Geofence, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, latitude, longitude, radius, created_at
		 FROM geofences ORDER BY created_at DESC`)

	if err != nil {
		return nil, fmt.Errorf("failed to query geofences: %w", err)
	}
	defer rows.Close()

	var geofences []domain.Geofence
	for rows.Next() {
		gf := domain.Geofence{}
		err := rows.Scan(&gf.ID, &gf.Name, &gf.Latitude, &gf.Longitude, &gf.Radius, &gf.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}
		geofences = append(geofences, gf)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating geofences: %w", err)
	}

	return geofences, nil
}

// CreateGeofence creates a new geofence
func (r *PostgresGeofenceRepository) CreateGeofence(ctx context.Context, geofence *domain.Geofence) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO geofences (id, name, latitude, longitude, radius, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		geofence.ID, geofence.Name, geofence.Latitude, geofence.Longitude, geofence.Radius, geofence.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create geofence: %w", err)
	}

	return nil
}

// GetGeofencesByPoint retrieves geofences that contain a point
func (r *PostgresGeofenceRepository) GetGeofencesByPoint(ctx context.Context, latitude float64, longitude float64) ([]domain.Geofence, error) {
	// Using PostGIS for geographic queries
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, latitude, longitude, radius, created_at
		 FROM geofences
		 WHERE ST_Distance(ST_Point(longitude, latitude), ST_Point($1, $2)) <= radius`,
		longitude, latitude)

	if err != nil {
		return nil, fmt.Errorf("failed to query geofences by point: %w", err)
	}
	defer rows.Close()

	var geofences []domain.Geofence
	for rows.Next() {
		gf := domain.Geofence{}
		err := rows.Scan(&gf.ID, &gf.Name, &gf.Latitude, &gf.Longitude, &gf.Radius, &gf.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}
		geofences = append(geofences, gf)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating geofences: %w", err)
	}

	return geofences, nil
}
