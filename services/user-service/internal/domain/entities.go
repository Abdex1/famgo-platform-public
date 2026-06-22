// services/user-service/internal/domain/entities.go
// User Service Domain Entities

package domain

import (
	"time"
)

// User represents a platform user (driver or passenger)
type User struct {
	ID        string    // UUID
	Phone     string    // Unique phone number
	Email     string    // Unique email
	FirstName string
	LastName  string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusDeleted   UserStatus = "DELETED"
)

// DriverProfile represents a driver's profile
type DriverProfile struct {
	ID                string    // UUID
	UserID            string    // Foreign key to user
	LicenseNumber     string    // Unique license
	LicenseExpiry     time.Time
	VehicleNumber     string    // License plate
	VehicleType       string    // Car, Motorcycle, etc.
	VerificationStatus VerificationStatus
	RatingCount       int32
	AverageRating     float32
	AcceptanceRate    float32 // Percentage
	CancellationRate  float32 // Percentage
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// PassengerProfile represents a passenger's profile
type PassengerProfile struct {
	ID                 string    // UUID
	UserID             string    // Foreign key to user
	PreferredLanguage  string
	EmergencyContact   string
	EmergencyPhone     string
	RatingCount        int32
	AverageRating      float32
	PreferredPayments  []string // Payment methods
	SavedLocations     map[string]SavedLocation
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SavedLocation represents a saved address
type SavedLocation struct {
	Label     string
	Latitude  float64
	Longitude float64
	Address   string
}

// VerificationStatus represents driver verification state
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "PENDING"
	VerificationStatusVerified VerificationStatus = "VERIFIED"
	VerificationStatusRejected VerificationStatus = "REJECTED"
	VerificationStatusExpired  VerificationStatus = "EXPIRED"
)

// UserPreference represents user preferences
type UserPreference struct {
	ID                   string    // UUID
	UserID               string    // Foreign key to user
	NotificationsEnabled bool
	ShareLocation        bool
	PreferredCurrency    string
	PreferredLanguage    string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// NewUserWithID creates a new user entity
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewUserWithID(id, phone, email, firstName, lastName string) *User {
	return &User{
		ID:        id,
		Phone:     phone,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Status:    UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewDriverProfileWithID creates a new driver profile
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewDriverProfileWithID(id, userID, licenseNumber, vehicleNumber, vehicleType string) *DriverProfile {
	return &DriverProfile{
		ID:                 id,
		UserID:             userID,
		LicenseNumber:      licenseNumber,
		VehicleNumber:      vehicleNumber,
		VehicleType:        vehicleType,
		VerificationStatus: VerificationStatusPending,
		RatingCount:        0,
		AverageRating:      0,
		AcceptanceRate:     100,
		CancellationRate:   0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// NewPassengerProfileWithID creates a new passenger profile
// ID is generated at application layer (Rule 4: domain has ZERO external dependencies)
func NewPassengerProfileWithID(id, userID string) *PassengerProfile {
	return &PassengerProfile{
		ID:                id,
		UserID:            userID,
		PreferredLanguage: "en",
		RatingCount:       0,
		AverageRating:     0,
		PreferredPayments: []string{"card", "cash"},
		SavedLocations:    make(map[string]SavedLocation),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// Activate marks user as active
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

// Suspend marks user as suspended
func (u *User) Suspend() {
	u.Status = UserStatusSuspended
	u.UpdatedAt = time.Now()
}

// Deactivate marks user as inactive
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

// UpdateDriverRating updates driver rating
func (dp *DriverProfile) UpdateRating(newRating float32, totalRatings int32) {
	// Calculate weighted average
	totalScore := dp.AverageRating * float32(dp.RatingCount)
	newTotalScore := totalScore + newRating
	dp.AverageRating = newTotalScore / float32(totalRatings)
	dp.RatingCount = totalRatings
	dp.UpdatedAt = time.Now()
}

// VerifyDriver marks driver as verified
func (dp *DriverProfile) Verify() {
	dp.VerificationStatus = VerificationStatusVerified
	dp.UpdatedAt = time.Now()
}

// RejectVerification marks verification as rejected
func (dp *DriverProfile) RejectVerification() {
	dp.VerificationStatus = VerificationStatusRejected
	dp.UpdatedAt = time.Now()
}
