// services/user-service/internal/domain/user_service.go
// User domain service with pure business logic

package domain

// UserService provides domain-level operations (pure logic, no I/O)
type UserService struct{}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// CanActivateUser validates if user can be activated
func (s *UserService) CanActivateUser(user *User) bool {
	return user.Status == UserStatusInactive || user.Status == UserStatusSuspended
}

// CanSuspendUser validates if user can be suspended
func (s *UserService) CanSuspendUser(user *User) bool {
	return user.Status == UserStatusActive
}

// CanDeleteUser validates if user can be deleted
func (s *UserService) CanDeleteUser(user *User) bool {
	return user.Status != UserStatusDeleted
}

// CanVerifyDriver validates if driver can be verified
func (s *UserService) CanVerifyDriver(profile *DriverProfile) bool {
	if profile.VerificationStatus == VerificationStatusVerified {
		return false // Already verified
	}
	if profile.VerificationStatus == VerificationStatusRejected {
		return true // Can reverify after rejection
	}
	return profile.VerificationStatus == VerificationStatusPending
}

// CalculateDriverAcceptanceRate calculates acceptance rate
func (s *UserService) CalculateDriverAcceptanceRate(accepted, total int32) float32 {
	if total == 0 {
		return 0
	}
	return (float32(accepted) / float32(total)) * 100
}

// CalculateDriverCancellationRate calculates cancellation rate
func (s *UserService) CalculateDriverCancellationRate(cancelled, total int32) float32 {
	if total == 0 {
		return 0
	}
	return (float32(cancelled) / float32(total)) * 100
}

// IsDriverVerified checks if driver is verified
func (s *UserService) IsDriverVerified(profile *DriverProfile) bool {
	return profile.VerificationStatus == VerificationStatusVerified
}

// IsLicenseExpired checks if license is expired
func (s *UserService) IsLicenseExpired(profile *DriverProfile) bool {
	// Import time at top
	return profile.LicenseExpiry.Before(getCurrentTime())
}

// ValidatePhoneNumber validates phone number format (basic)
func (s *UserService) ValidatePhoneNumber(phone string) bool {
	// Remove spaces and special chars
	cleaned := ""
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' || ch == '+' {
			cleaned += string(ch)
		}
	}
	// Must be at least 7 digits
	return len(cleaned) >= 7 && len(cleaned) <= 15
}

// ValidateEmail validates email format (basic)
func (s *UserService) ValidateEmail(email string) bool {
	// Basic email validation
	parts := split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !contains(parts[1], ".") {
		return false
	}
	return true
}

// Helper functions
func getCurrentTime() time.Time {
	return time.Now()
}

func split(s, sep string) []string {
	var result []string
	var current string
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

import "time"
