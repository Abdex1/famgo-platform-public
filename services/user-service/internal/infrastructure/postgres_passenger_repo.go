// services/user-service/internal/infrastructure/postgres_passenger_repo.go
// PostgreSQL Passenger Profile Repository

package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// PostgresPassengerProfileRepository implements PassengerProfileRepository
type PostgresPassengerProfileRepository struct {
	db *sql.DB
}

func NewPostgresPassengerProfileRepository(db *sql.DB) *PostgresPassengerProfileRepository {
	return &PostgresPassengerProfileRepository{db: db}
}

func (r *PostgresPassengerProfileRepository) GetProfile(ctx context.Context, profileID string) (*domain.PassengerProfile, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, preferred_language, emergency_contact, emergency_phone,
                rating_count, average_rating, preferred_payments, saved_locations, created_at, updated_at
         FROM passenger_profiles WHERE id = $1`,
		profileID)

	profile := &domain.PassengerProfile{}
	var paymentsJSON, locationsJSON sql.NullString

	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.PreferredLanguage,
		&profile.EmergencyContact,
		&profile.EmergencyPhone,
		&profile.RatingCount,
		&profile.AverageRating,
		&paymentsJSON,
		&locationsJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal JSON
	if paymentsJSON.Valid {
		json.Unmarshal([]byte(paymentsJSON.String), &profile.PreferredPayments)
	}
	if locationsJSON.Valid {
		json.Unmarshal([]byte(locationsJSON.String), &profile.SavedLocations)
	}

	return profile, nil
}

func (r *PostgresPassengerProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.PassengerProfile, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, preferred_language, emergency_contact, emergency_phone,
                rating_count, average_rating, preferred_payments, saved_locations, created_at, updated_at
         FROM passenger_profiles WHERE user_id = $1`,
		userID)

	profile := &domain.PassengerProfile{}
	var paymentsJSON, locationsJSON sql.NullString

	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.PreferredLanguage,
		&profile.EmergencyContact,
		&profile.EmergencyPhone,
		&profile.RatingCount,
		&profile.AverageRating,
		&paymentsJSON,
		&locationsJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if paymentsJSON.Valid {
		json.Unmarshal([]byte(paymentsJSON.String), &profile.PreferredPayments)
	}
	if locationsJSON.Valid {
		json.Unmarshal([]byte(locationsJSON.String), &profile.SavedLocations)
	}

	return profile, nil
}

func (r *PostgresPassengerProfileRepository) CreateProfile(ctx context.Context, profile *domain.PassengerProfile) error {
	paymentsJSON, _ := json.Marshal(profile.PreferredPayments)
	locationsJSON, _ := json.Marshal(profile.SavedLocations)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO passenger_profiles (id, user_id, preferred_language, emergency_contact, emergency_phone,
                rating_count, average_rating, preferred_payments, saved_locations, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		profile.ID,
		profile.UserID,
		profile.PreferredLanguage,
		profile.EmergencyContact,
		profile.EmergencyPhone,
		profile.RatingCount,
		profile.AverageRating,
		string(paymentsJSON),
		string(locationsJSON),
		profile.CreatedAt,
		profile.UpdatedAt,
	)
	return err
}

func (r *PostgresPassengerProfileRepository) UpdateProfile(ctx context.Context, profile *domain.PassengerProfile) error {
	paymentsJSON, _ := json.Marshal(profile.PreferredPayments)
	locationsJSON, _ := json.Marshal(profile.SavedLocations)

	_, err := r.db.ExecContext(ctx,
		`UPDATE passenger_profiles SET preferred_language = $1, emergency_contact = $2, emergency_phone = $3,
                preferred_payments = $4, saved_locations = $5, updated_at = $6 WHERE id = $7`,
		profile.PreferredLanguage,
		profile.EmergencyContact,
		profile.EmergencyPhone,
		string(paymentsJSON),
		string(locationsJSON),
		profile.UpdatedAt,
		profile.ID,
	)
	return err
}

func (r *PostgresPassengerProfileRepository) UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error {
	row := r.db.QueryRowContext(ctx,
		`SELECT average_rating, rating_count FROM passenger_profiles WHERE id = $1`,
		profileID)

	var currentRating float32
	var currentCount int32
	if err := row.Scan(&currentRating, &currentCount); err != nil {
		return err
	}

	totalScore := currentRating*float32(currentCount) + newRating
	newAverage := totalScore / float32(totalRatings)

	_, err := r.db.ExecContext(ctx,
		`UPDATE passenger_profiles SET average_rating = $1, rating_count = $2, updated_at = NOW() WHERE id = $3`,
		newAverage, totalRatings, profileID)

	return err
}

func (r *PostgresPassengerProfileRepository) AddSavedLocation(ctx context.Context, profileID string, label string, location domain.SavedLocation) error {
	// Get current profile
	profile, err := r.GetProfile(ctx, profileID)
	if err != nil {
		return err
	}

	// Add location
	profile.SavedLocations[label] = location

	// Update
	locationsJSON, _ := json.Marshal(profile.SavedLocations)
	_, err = r.db.ExecContext(ctx,
		`UPDATE passenger_profiles SET saved_locations = $1, updated_at = NOW() WHERE id = $2`,
		string(locationsJSON), profileID)

	return err
}

func (r *PostgresPassengerProfileRepository) RemoveSavedLocation(ctx context.Context, profileID string, label string) error {
	// Get current profile
	profile, err := r.GetProfile(ctx, profileID)
	if err != nil {
		return err
	}

	// Remove location
	delete(profile.SavedLocations, label)

	// Update
	locationsJSON, _ := json.Marshal(profile.SavedLocations)
	_, err = r.db.ExecContext(ctx,
		`UPDATE passenger_profiles SET saved_locations = $1, updated_at = NOW() WHERE id = $2`,
		string(locationsJSON), profileID)

	return err
}
