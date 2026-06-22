// services/user-service/internal/application/interfaces.go
// Application layer interfaces (what application depends on)

package application

import (
	"context"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// UserRepository interface
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

// DriverProfileRepository interface
type DriverProfileRepository interface {
	GetProfile(ctx context.Context, profileID string) (*domain.DriverProfile, error)
	GetByUserID(ctx context.Context, userID string) (*domain.DriverProfile, error)
	CreateProfile(ctx context.Context, profile *domain.DriverProfile) error
	UpdateProfile(ctx context.Context, profile *domain.DriverProfile) error
	UpdateVerification(ctx context.Context, profileID string, status domain.VerificationStatus) error
	UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error
	ListActiveDrivers(ctx context.Context, limit, offset int) ([]domain.DriverProfile, error)
}

// PassengerProfileRepository interface
type PassengerProfileRepository interface {
	GetProfile(ctx context.Context, profileID string) (*domain.PassengerProfile, error)
	GetByUserID(ctx context.Context, userID string) (*domain.PassengerProfile, error)
	CreateProfile(ctx context.Context, profile *domain.PassengerProfile) error
	UpdateProfile(ctx context.Context, profile *domain.PassengerProfile) error
	UpdateRating(ctx context.Context, profileID string, newRating float32, totalRatings int32) error
	AddSavedLocation(ctx context.Context, profileID string, label string, location domain.SavedLocation) error
	RemoveSavedLocation(ctx context.Context, profileID string, label string) error
}

// UserCache interface
type UserCache interface {
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	SetUser(ctx context.Context, user *domain.User, ttl int32) error
	DeleteUser(ctx context.Context, userID string) error
	GetDriverProfile(ctx context.Context, driverID string) (*domain.DriverProfile, error)
	SetDriverProfile(ctx context.Context, profile *domain.DriverProfile, ttl int32) error
	DeleteDriverProfile(ctx context.Context, driverID string) error
	GetPassengerProfile(ctx context.Context, passengerID string) (*domain.PassengerProfile, error)
	SetPassengerProfile(ctx context.Context, profile *domain.PassengerProfile, ttl int32) error
	DeletePassengerProfile(ctx context.Context, passengerID string) error
}
