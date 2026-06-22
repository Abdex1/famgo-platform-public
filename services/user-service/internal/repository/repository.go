package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"famgo/user-service/internal/model"
)

// UserRepository handles database operations for user profiles
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateProfile creates a new user profile
func (r *UserRepository) CreateProfile(ctx context.Context, profile *model.UserProfile) error {
	query := `
	INSERT INTO user_profiles (id, auth_id, first_name, last_name, profile_picture_url, email_verified, phone_verified, rating, total_rides, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, auth_id, first_name, last_name, profile_picture_url, email_verified, phone_verified, rating, total_rides, created_at, updated_at
	`

	return r.db.QueryRowxContext(ctx, query,
		uuid.New().String(),
		profile.AuthID,
		profile.FirstName,
		profile.LastName,
		profile.ProfilePictureURL,
		profile.EmailVerified,
		profile.PhoneVerified,
		profile.Rating,
		profile.TotalRides,
		time.Now(),
		time.Now(),
	).StructScan(profile)
}

// GetProfileByAuthID retrieves profile by auth ID
func (r *UserRepository) GetProfileByAuthID(ctx context.Context, authID string) (*model.UserProfile, error) {
	query := `
	SELECT id, auth_id, first_name, last_name, profile_picture_url, email_verified, phone_verified, rating, total_rides, created_at, updated_at
	FROM user_profiles
	WHERE auth_id = $1
	`

	profile := &model.UserProfile{}
	if err := r.db.QueryRowxContext(ctx, query, authID).StructScan(profile); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("profile not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return profile, nil
}

// GetProfileByID retrieves profile by ID
func (r *UserRepository) GetProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	query := `
	SELECT id, auth_id, first_name, last_name, profile_picture_url, email_verified, phone_verified, rating, total_rides, created_at, updated_at
	FROM user_profiles
	WHERE id = $1
	`

	profile := &model.UserProfile{}
	if err := r.db.QueryRowxContext(ctx, query, id).StructScan(profile); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("profile not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return profile, nil
}

// UpdateProfile updates user profile
func (r *UserRepository) UpdateProfile(ctx context.Context, profile *model.UserProfile) error {
	query := `
	UPDATE user_profiles
	SET first_name = $1, last_name = $2, profile_picture_url = $3, updated_at = $4
	WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		profile.FirstName,
		profile.LastName,
		profile.ProfilePictureURL,
		time.Now(),
		profile.ID,
	)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return errors.New("profile not found")
	}

	return nil
}

// PreferencesRepository handles database operations for preferences
type PreferencesRepository struct {
	db *sqlx.DB
}

// NewPreferencesRepository creates a new preferences repository
func NewPreferencesRepository(db *sqlx.DB) *PreferencesRepository {
	return &PreferencesRepository{db: db}
}

// CreatePreferences creates default preferences for user
func (r *PreferencesRepository) CreatePreferences(ctx context.Context, userID string) (*model.UserPreferences, error) {
	prefs := &model.UserPreferences{
		ID:               uuid.New().String(),
		UserID:           userID,
		NotificationEmail: true,
		NotificationSMS:  true,
		Language:         "en",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `
	INSERT INTO user_preferences (id, user_id, notification_email, notification_sms, language, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, user_id, notification_email, notification_sms, language, created_at, updated_at
	`

	return prefs, r.db.QueryRowxContext(ctx, query,
		prefs.ID, prefs.UserID, prefs.NotificationEmail, prefs.NotificationSMS, prefs.Language, prefs.CreatedAt, prefs.UpdatedAt,
	).StructScan(prefs)
}

// GetPreferencesByUserID retrieves preferences by user ID
func (r *PreferencesRepository) GetPreferencesByUserID(ctx context.Context, userID string) (*model.UserPreferences, error) {
	query := `
	SELECT id, user_id, notification_email, notification_sms, language, created_at, updated_at
	FROM user_preferences
	WHERE user_id = $1
	`

	prefs := &model.UserPreferences{}
	if err := r.db.QueryRowxContext(ctx, query, userID).StructScan(prefs); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("preferences not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return prefs, nil
}

// UpdatePreferences updates user preferences
func (r *PreferencesRepository) UpdatePreferences(ctx context.Context, prefs *model.UserPreferences) error {
	query := `
	UPDATE user_preferences
	SET notification_email = $1, notification_sms = $2, language = $3, updated_at = $4
	WHERE user_id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		prefs.NotificationEmail,
		prefs.NotificationSMS,
		prefs.Language,
		time.Now(),
		prefs.UserID,
	)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return errors.New("preferences not found")
	}

	return nil
}

// AddressRepository handles database operations for addresses
type AddressRepository struct {
	db *sqlx.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *sqlx.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// CreateAddress creates a new saved address
func (r *AddressRepository) CreateAddress(ctx context.Context, userID string, addr *model.AddressRequest) (*model.UserAddress, error) {
	address := &model.UserAddress{
		ID:           uuid.New().String(),
		UserID:       userID,
		Type:         addr.Type,
		AddressLine1: addr.AddressLine1,
		City:         addr.City,
		Lat:          addr.Lat,
		Lng:          addr.Lng,
		CreatedAt:    time.Now(),
	}

	query := `
	INSERT INTO user_addresses (id, user_id, type, address_line_1, city, lat, lng, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, user_id, type, address_line_1, city, lat, lng, created_at
	`

	return address, r.db.QueryRowxContext(ctx, query,
		address.ID, address.UserID, address.Type, address.AddressLine1, address.City, address.Lat, address.Lng, address.CreatedAt,
	).StructScan(address)
}

// GetAddressesByUserID retrieves all addresses for user
func (r *AddressRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]*model.UserAddress, error) {
	query := `
	SELECT id, user_id, type, address_line_1, city, lat, lng, created_at
	FROM user_addresses
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	var addresses []*model.UserAddress
	if err := r.db.SelectContext(ctx, &addresses, query, userID); err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return addresses, nil
}

// DeleteAddress deletes a saved address
func (r *AddressRepository) DeleteAddress(ctx context.Context, addressID string) error {
	query := `DELETE FROM user_addresses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, addressID)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return errors.New("address not found")
	}

	return nil
}
