// services/user-service/internal/infrastructure/postgres_driver_repo.go
// PostgreSQL Driver Profile Repository

package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// PostgresDriverProfileRepository implements DriverProfileRepository
type PostgresDriverProfileRepository struct {
	db *sql.DB
}

func NewPostgresDriverProfileRepository(db *sql.DB) *PostgresDriverProfileRepository {
	return &PostgresDriverProfileRepository{db: db}
}

func (r *PostgresDriverProfileRepository) GetProfile(ctx context.Context, profileID string) (*domain.DriverProfile, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, license_number, license_expiry, vehicle_number, vehicle_type,
                verification_status, rating_count, average_rating, acceptance_rate, cancellation_rate,
                created_at, updated_at FROM driver_profiles WHERE id = $1`,
		profileID)

	profile := &domain.DriverProfile{}
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.LicenseNumber,
		&profile.LicenseExpiry,
		&profile.VehicleNumber,
		&profile.VehicleType,
		&profile.VerificationStatus,
		&profile.RatingCount,
		&profile.AverageRating,
		&profile.AcceptanceRate,
		&profile.CancellationRate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func (r *PostgresDriverProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.DriverProfile, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, license_number, license_expiry, vehicle_number, vehicle_type,
                verification_status, rating_count, average_rating, acceptance_rate, cancellation_rate,
                created_at, updated_at FROM driver_profiles WHERE user_id = $1`,
		userID)

	profile := &domain.DriverProfile{}
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.LicenseNumber,
		&profile.LicenseExpiry,
		&profile.VehicleNumber,
		&profile.VehicleType,
		&profile.VerificationStatus,
		&profile.RatingCount,
		&profile.AverageRating,
		&profile.AcceptanceRate,
		&profile.CancellationRate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func (r *PostgresDriverProfileRepository) CreateProfile(ctx context.Context, profile *domain.DriverProfile) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO driver_profiles (id, user_id, license_number, license_expiry, vehicle_number, vehicle_type,
                verification_status, rating_count, average_rating, acceptance_rate, cancellation_rate, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		profile.ID,
		profile.UserID,
		profile.LicenseNumber,
		profile.LicenseExpiry,
		profile.VehicleNumber,
		profile.VehicleType,
		profile.VerificationStatus,
		profile.RatingCount,
		profile.AverageRating,
		profile.AcceptanceRate,
		profile.CancellationRate,
		profile.CreatedAt,
		profile.UpdatedAt,
	)
	return err
}

func (r *PostgresDriverProfileRepository) UpdateProfile(ctx context.Context, profile *domain.DriverProfile) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE driver_profiles SET license_number = $1, license_expiry = $2, vehicle_number = $3, 
         vehicle_type = $4, verification_status = $5, rating_count = $6, average_rating = $7,
         acceptance_rate = $8, cancellation_rate = $9, updated_at = $10 WHERE id = $11`,
		profile.LicenseNumber,
		profile.LicenseExpiry,
		profile.VehicleNumber,
		profile.VehicleType,
		profile.VerificationStatus,
		profile.RatingCount,
		profile.AverageRating,
		profile.AcceptanceRate,
		profile.CancellationRate,
		profile.UpdatedAt,
		profile.ID,
	)
	return err
}

func (r *PostgresDriverProfileRepository) UpdateVerification(ctx context.Context, profileID string, status domain.VerificationStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE driver_profiles SET verification_status = $1, updated_at = NOW() WHERE id = $2`,
		status, profileID)
	return err
}

func (r *PostgresDriverProfileRepository) UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error {
	// Get current rating
	row := r.db.QueryRowContext(ctx,
		`SELECT average_rating, rating_count FROM driver_profiles WHERE id = $1`,
		profileID)

	var currentRating float32
	var currentCount int32
	if err := row.Scan(&currentRating, &currentCount); err != nil {
		return err
	}

	// Calculate new average
	totalScore := currentRating*float32(currentCount) + newRating
	newAverage := totalScore / float32(totalRatings)

	// Update
	_, err := r.db.ExecContext(ctx,
		`UPDATE driver_profiles SET average_rating = $1, rating_count = $2, updated_at = NOW() WHERE id = $3`,
		newAverage, totalRatings, profileID)

	return err
}

func (r *PostgresDriverProfileRepository) ListActiveDrivers(ctx context.Context, limit, offset int) ([]domain.DriverProfile, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, license_number, license_expiry, vehicle_number, vehicle_type,
                verification_status, rating_count, average_rating, acceptance_rate, cancellation_rate,
                created_at, updated_at FROM driver_profiles 
         WHERE verification_status = 'VERIFIED' LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []domain.DriverProfile
	for rows.Next() {
		profile := domain.DriverProfile{}
		err := rows.Scan(
			&profile.ID,
			&profile.UserID,
			&profile.LicenseNumber,
			&profile.LicenseExpiry,
			&profile.VehicleNumber,
			&profile.VehicleType,
			&profile.VerificationStatus,
			&profile.RatingCount,
			&profile.AverageRating,
			&profile.AcceptanceRate,
			&profile.CancellationRate,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, rows.Err()
}
