// services/user-service/tests/unit/user_service_test.go
// User Service Domain Logic Tests

package unit

import (
	"testing"
	"time"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// TestValidateEmail tests email validation
func TestValidateEmail(t *testing.T) {
	service := domain.NewUserService()

	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"test.email@domain.co.uk", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := service.ValidateEmail(tt.email)
			if result != tt.valid {
				t.Errorf("ValidateEmail(%s) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}

// TestValidatePhoneNumber tests phone validation
func TestValidatePhoneNumber(t *testing.T) {
	service := domain.NewUserService()

	tests := []struct {
		phone string
		valid bool
	}{
		{"+1234567890", true},
		{"1234567890", true},
		{"+1 (234) 567-8900", true},
		{"123456", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			result := service.ValidatePhoneNumber(tt.phone)
			if result != tt.valid {
				t.Errorf("ValidatePhoneNumber(%s) = %v, want %v", tt.phone, result, tt.valid)
			}
		})
	}
}

// TestCanActivateUser tests user activation logic
func TestCanActivateUser(t *testing.T) {
	service := domain.NewUserService()

	tests := []struct {
		name   string
		status domain.UserStatus
		can    bool
	}{
		{"active user", domain.UserStatusActive, false},
		{"inactive user", domain.UserStatusInactive, true},
		{"suspended user", domain.UserStatusSuspended, true},
		{"deleted user", domain.UserStatusDeleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{Status: tt.status}
			result := service.CanActivateUser(user)
			if result != tt.can {
				t.Errorf("CanActivateUser(%s) = %v, want %v", tt.status, result, tt.can)
			}
		})
	}
}

// TestCanVerifyDriver tests driver verification logic
func TestCanVerifyDriver(t *testing.T) {
	service := domain.NewUserService()

	tests := []struct {
		name   string
		status domain.VerificationStatus
		can    bool
	}{
		{"pending driver", domain.VerificationStatusPending, true},
		{"verified driver", domain.VerificationStatusVerified, false},
		{"rejected driver", domain.VerificationStatusRejected, true},
		{"expired driver", domain.VerificationStatusExpired, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &domain.DriverProfile{VerificationStatus: tt.status}
			result := service.CanVerifyDriver(profile)
			if result != tt.can {
				t.Errorf("CanVerifyDriver(%s) = %v, want %v", tt.status, result, tt.can)
			}
		})
	}
}

// TestCalculateDriverAcceptanceRate tests acceptance rate calculation
func TestCalculateDriverAcceptanceRate(t *testing.T) {
	service := domain.NewUserService()

	tests := []struct {
		accepted int32
		total    int32
		expected float32
	}{
		{100, 100, 100},
		{50, 100, 50},
		{0, 100, 0},
		{0, 0, 0},
		{95, 100, 95},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := service.CalculateDriverAcceptanceRate(tt.accepted, tt.total)
			if result != tt.expected {
				t.Errorf("CalculateDriverAcceptanceRate(%d, %d) = %v, want %v", tt.accepted, tt.total, result, tt.expected)
			}
		})
	}
}

// TestNewUser tests user creation
func TestNewUser(t *testing.T) {
	phone := "+1234567890"
	email := "user@example.com"
	firstName := "John"
	lastName := "Doe"

	user := domain.NewUser(phone, email, firstName, lastName)

	if user.Phone != phone {
		t.Errorf("Phone = %v, want %v", user.Phone, phone)
	}
	if user.Email != email {
		t.Errorf("Email = %v, want %v", user.Email, email)
	}
	if user.FirstName != firstName {
		t.Errorf("FirstName = %v, want %v", user.FirstName, firstName)
	}
	if user.Status != domain.UserStatusActive {
		t.Errorf("Status = %v, want %v", user.Status, domain.UserStatusActive)
	}
	if user.ID == "" {
		t.Error("ID should not be empty")
	}
}

// TestNewDriverProfile tests driver profile creation
func TestNewDriverProfile(t *testing.T) {
	userID := "user-123"
	licenseNumber := "DL123456"
	vehicleNumber := "ABC123"
	vehicleType := "sedan"

	profile := domain.NewDriverProfile(userID, licenseNumber, vehicleNumber, vehicleType)

	if profile.UserID != userID {
		t.Errorf("UserID = %v, want %v", profile.UserID, userID)
	}
	if profile.LicenseNumber != licenseNumber {
		t.Errorf("LicenseNumber = %v, want %v", profile.LicenseNumber, licenseNumber)
	}
	if profile.VerificationStatus != domain.VerificationStatusPending {
		t.Errorf("VerificationStatus = %v, want %v", profile.VerificationStatus, domain.VerificationStatusPending)
	}
	if profile.AcceptanceRate != 100 {
		t.Errorf("AcceptanceRate = %v, want %v", profile.AcceptanceRate, 100)
	}
}

// TestUserActivation tests user activation state transition
func TestUserActivation(t *testing.T) {
	user := domain.NewUser("+1234567890", "user@example.com", "John", "Doe")
	user.Status = domain.UserStatusInactive

	user.Activate()

	if user.Status != domain.UserStatusActive {
		t.Errorf("Status = %v, want %v", user.Status, domain.UserStatusActive)
	}
}

// TestDriverVerification tests driver verification state transition
func TestDriverVerification(t *testing.T) {
	profile := domain.NewDriverProfile("user-123", "DL123456", "ABC123", "sedan")

	profile.Verify()

	if profile.VerificationStatus != domain.VerificationStatusVerified {
		t.Errorf("VerificationStatus = %v, want %v", profile.VerificationStatus, domain.VerificationStatusVerified)
	}
}

// TestUpdateDriverRating tests rating update
func TestUpdateDriverRating(t *testing.T) {
	profile := domain.NewDriverProfile("user-123", "DL123456", "ABC123", "sedan")

	// First rating
	profile.UpdateRating(5.0, 1)
	if profile.AverageRating != 5.0 {
		t.Errorf("AverageRating = %v, want %v", profile.AverageRating, 5.0)
	}
	if profile.RatingCount != 1 {
		t.Errorf("RatingCount = %v, want %v", profile.RatingCount, 1)
	}

	// Second rating
	profile.UpdateRating(4.0, 2)
	expectedRating := float32(4.5)
	if profile.AverageRating != expectedRating {
		t.Errorf("AverageRating = %v, want %v", profile.AverageRating, expectedRating)
	}
	if profile.RatingCount != 2 {
		t.Errorf("RatingCount = %v, want %v", profile.RatingCount, 2)
	}
}
