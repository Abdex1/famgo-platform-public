package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"famgo/driver-service/internal/model"
)

// DriverRepository handles database operations for drivers
type DriverRepository struct {
	db *sqlx.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *sqlx.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// DB returns the underlying database handle.
func (r *DriverRepository) DB() *sqlx.DB {
	return r.db
}

// CreateDriver creates a new driver
func (r *DriverRepository) CreateDriver(ctx context.Context, driver *model.Driver) error {
	driver.ID = uuid.New().String()
	driver.Status = model.DriverStatePending
	driver.VerificationStatus = "pending"
	driver.DateJoined = time.Now()
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()

	query := `
	INSERT INTO drivers (id, auth_id, license_number, license_expiry, status, verification_status, date_joined, rating, total_rides, total_earnings, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, auth_id, license_number, license_expiry, status, verification_status, date_joined, rating, total_rides, total_earnings, created_at, updated_at
	`

	return r.db.QueryRowxContext(ctx, query,
		driver.ID,
		driver.AuthID,
		driver.LicenseNumber,
		driver.LicenseExpiry,
		driver.Status,
		driver.VerificationStatus,
		driver.DateJoined,
		driver.Rating,
		driver.TotalRides,
		driver.TotalEarnings,
		driver.CreatedAt,
		driver.UpdatedAt,
	).StructScan(driver)
}

// GetDriverByAuthID retrieves driver by auth ID
func (r *DriverRepository) GetDriverByAuthID(ctx context.Context, authID string) (*model.Driver, error) {
	query := `
	SELECT id, auth_id, license_number, license_expiry, status, verification_status, date_joined, rating, total_rides, total_earnings, created_at, updated_at
	FROM drivers
	WHERE auth_id = $1
	`

	driver := &model.Driver{}
	if err := r.db.QueryRowxContext(ctx, query, authID).StructScan(driver); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("driver not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return driver, nil
}

// GetDriverByID retrieves driver by ID
func (r *DriverRepository) GetDriverByID(ctx context.Context, id string) (*model.Driver, error) {
	query := `
	SELECT id, auth_id, license_number, license_expiry, status, verification_status, date_joined, rating, total_rides, total_earnings, created_at, updated_at
	FROM drivers
	WHERE id = $1
	`

	driver := &model.Driver{}
	if err := r.db.QueryRowxContext(ctx, query, id).StructScan(driver); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("driver not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return driver, nil
}

// UpdateDriver updates driver information
func (r *DriverRepository) UpdateDriver(ctx context.Context, driver *model.Driver) error {
	driver.UpdatedAt = time.Now()

	query := `
	UPDATE drivers
	SET license_number = $1, license_expiry = $2, status = $3, verification_status = $4, rating = $5, total_rides = $6, total_earnings = $7, updated_at = $8
	WHERE id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		driver.LicenseNumber,
		driver.LicenseExpiry,
		driver.Status,
		driver.VerificationStatus,
		driver.Rating,
		driver.TotalRides,
		driver.TotalEarnings,
		driver.UpdatedAt,
		driver.ID,
	)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return errors.New("driver not found")
	}

	return nil
}

// DriverStateRepository handles state machine operations (Pattern 4)
type DriverStateRepository struct {
	db *sqlx.DB
}

// NewDriverStateRepository creates a new driver state repository
func NewDriverStateRepository(db *sqlx.DB) *DriverStateRepository {
	return &DriverStateRepository{db: db}
}

// TransitionState transitions driver to a new state (Pattern 4: State Machine)
func (r *DriverStateRepository) TransitionState(ctx context.Context, driverID string, fromState, toState, reason string) (*model.DriverState, error) {
	// Validate transition
	if !model.IsValidTransition(fromState, toState) {
		return nil, fmt.Errorf("invalid state transition: %s -> %s", fromState, toState)
	}

	state := &model.DriverState{
		ID:            uuid.New().String(),
		DriverID:      driverID,
		CurrentState:  toState,
		PreviousState: fromState,
		Reason:        reason,
		TransitionAt:  time.Now(),
		CreatedAt:     time.Now(),
	}

	query := `
	INSERT INTO driver_states (id, driver_id, current_state, previous_state, reason, transition_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		state.ID,
		state.DriverID,
		state.CurrentState,
		state.PreviousState,
		state.Reason,
		state.TransitionAt,
		state.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to record state transition: %w", err)
	}

	return state, nil
}

// GetCurrentState retrieves the current state of a driver
func (r *DriverStateRepository) GetCurrentState(ctx context.Context, driverID string) (*model.DriverState, error) {
	query := `
	SELECT id, driver_id, current_state, previous_state, reason, transition_at, created_at
	FROM driver_states
	WHERE driver_id = $1
	ORDER BY transition_at DESC
	LIMIT 1
	`

	state := &model.DriverState{}
	if err := r.db.QueryRowxContext(ctx, query, driverID).StructScan(state); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("no state history found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return state, nil
}

// GetStateHistory retrieves state transition history
func (r *DriverStateRepository) GetStateHistory(ctx context.Context, driverID string, limit int) ([]*model.DriverState, error) {
	query := `
	SELECT id, driver_id, current_state, previous_state, reason, transition_at, created_at
	FROM driver_states
	WHERE driver_id = $1
	ORDER BY transition_at DESC
	LIMIT $2
	`

	var states []*model.DriverState
	if err := r.db.SelectContext(ctx, &states, query, driverID, limit); err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return states, nil
}
