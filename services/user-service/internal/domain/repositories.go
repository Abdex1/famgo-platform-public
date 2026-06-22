// services/user-service/internal/domain/repositories.go
// Repository interfaces - what application depends on

package domain

import "context"

// UserRepository defines operations on users
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)
	DeleteUser(ctx context.Context, userID string) error
}

// DriverProfileRepository defines operations on driver profiles
type DriverProfileRepository interface {
	GetProfile(ctx context.Context, profileID string) (*DriverProfile, error)
	GetByUserID(ctx context.Context, userID string) (*DriverProfile, error)
	CreateProfile(ctx context.Context, profile *DriverProfile) error
	UpdateProfile(ctx context.Context, profile *DriverProfile) error
	UpdateVerification(ctx context.Context, profileID string, status VerificationStatus) error
	UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error
	ListActiveDrivers(ctx context.Context, limit, offset int) ([]DriverProfile, error)
}

// PassengerProfileRepository defines operations on passenger profiles
type PassengerProfileRepository interface {
	GetProfile(ctx context.Context, profileID string) (*PassengerProfile, error)
	GetByUserID(ctx context.Context, userID string) (*PassengerProfile, error)
	CreateProfile(ctx context.Context, profile *PassengerProfile) error
	UpdateProfile(ctx context.Context, profile *PassengerProfile) error
	UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error
	AddSavedLocation(ctx context.Context, profileID string, label string, location SavedLocation) error
	RemoveSavedLocation(ctx context.Context, profileID string, label string) error
}

// UserPreferenceRepository defines operations on user preferences
type UserPreferenceRepository interface {
	GetPreferences(ctx context.Context, userID string) (*UserPreference, error)
	CreatePreferences(ctx context.Context, pref *UserPreference) error
	UpdatePreferences(ctx context.Context, pref *UserPreference) error
}

// UserCache defines caching operations for users
type UserCache interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	SetUser(ctx context.Context, user *User, ttl int32) error
	DeleteUser(ctx context.Context, userID string) error
	GetDriverProfile(ctx context.Context, driverID string) (*DriverProfile, error)
	SetDriverProfile(ctx context.Context, profile *DriverProfile, ttl int32) error
	DeleteDriverProfile(ctx context.Context, driverID string) error
	GetPassengerProfile(ctx context.Context, passengerID string) (*PassengerProfile, error)
	SetPassengerProfile(ctx context.Context, profile *PassengerProfile, ttl int32) error
	DeletePassengerProfile(ctx context.Context, passengerID string) error
}
