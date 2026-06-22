package model

import "time"

// Driver represents a driver profile
type Driver struct {
	ID                 string    `db:"id"`
	AuthID             string    `db:"auth_id"`
	LicenseNumber      string    `db:"license_number"`
	LicenseExpiry      *time.Time `db:"license_expiry"`
	Status             string    `db:"status"`              // pending, approved, active, suspended
	VerificationStatus string    `db:"verification_status"` // pending, in_progress, approved, rejected
	DateJoined         time.Time `db:"date_joined"`
	Rating             float64   `db:"rating"`
	TotalRides         int       `db:"total_rides"`
	TotalEarnings      float64   `db:"total_earnings"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// DriverState represents a driver state transition (Pattern 4: State Machine)
type DriverState struct {
	ID             string    `db:"id"`
	DriverID       string    `db:"driver_id"`
	CurrentState   string    `db:"current_state"`   // pending, approved, active, suspended
	PreviousState  string    `db:"previous_state"`
	Reason         string    `db:"reason"`
	TransitionAt   time.Time `db:"transition_at"`
	CreatedAt      time.Time `db:"created_at"`
}

// Valid driver states (Pattern 4: State Machine)
const (
	DriverStatePending   = "pending"   // Initial state after registration
	DriverStateApproved  = "approved"  // After verification (WEEK 3)
	DriverStateActive    = "active"    // Ready to accept rides
	DriverStateSuspended = "suspended" // Suspended from platform
	DriverStateInactive  = "inactive"  // Offline
)

// Valid state transitions (Pattern 4)
var validTransitions = map[string][]string{
	DriverStatePending:   {DriverStateApproved, DriverStateSuspended},
	DriverStateApproved:  {DriverStateActive, DriverStateSuspended},
	DriverStateActive:    {DriverStateInactive, DriverStateSuspended},
	DriverStateInactive:  {DriverStateActive, DriverStateSuspended},
	DriverStateSuspended: {DriverStatePending}, // Can appeal
}

// IsValidTransition checks if transition is allowed
func IsValidTransition(from, to string) bool {
	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}
	for _, state := range allowedStates {
		if state == to {
			return true
		}
	}
	return false
}

// RegistrationRequest represents driver registration (step 1)
type RegistrationRequest struct {
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	LicenseNumber string `json:"license_number"`
	LicenseExpiry string `json:"license_expiry"` // YYYY-MM-DD
}

// VerifyRegistrationRequest represents verification (step 2)
type VerifyRegistrationRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// GetProfileResponse represents driver profile response
type GetProfileResponse struct {
	ID                 string  `json:"id"`
	Email              string  `json:"email"`
	Phone              string  `json:"phone"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	Status             string  `json:"status"`
	VerificationStatus string  `json:"verification_status"`
	Rating             float64 `json:"rating"`
	TotalRides         int     `json:"total_rides"`
	TotalEarnings      float64 `json:"total_earnings"`
}

// StateTransitionRequest represents a state transition request
type StateTransitionRequest struct {
	NewState string `json:"new_state"`
	Reason   string `json:"reason"`
}
