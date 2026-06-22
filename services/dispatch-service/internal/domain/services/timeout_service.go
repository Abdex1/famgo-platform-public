package services

import (
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
)

// TimeoutService enforces matching session expiry and timeout transitions.
type TimeoutService struct {
	defaultTTL time.Duration
}

func NewTimeoutService(defaultTTL time.Duration) *TimeoutService {
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}
	return &TimeoutService{defaultTTL: defaultTTL}
}

// ApplyExpiry sets the request expiry timestamp when not already configured.
func (s *TimeoutService) ApplyExpiry(request *entities.DispatchRequest) {
	if request.ExpiryTime.IsZero() {
		request.ExpiryTime = time.Now().UTC().Add(s.defaultTTL)
	}
}

// CheckAndExpire marks expired requests that are still active.
func (s *TimeoutService) CheckAndExpire(request *entities.DispatchRequest) error {
	if request == nil {
		return fmt.Errorf("dispatch request is required")
	}
	if !request.IsExpired() {
		return nil
	}
	if request.Status == entities.StatusExpired ||
		request.Status == entities.StatusCompleted ||
		request.Status == entities.StatusCancelled {
		return nil
	}
	return request.Expire()
}

// ShouldRetryAfterTimeout indicates whether matching can continue after timeout rejection.
func (s *TimeoutService) ShouldRetryAfterTimeout(request *entities.DispatchRequest) bool {
	return request != nil && request.CanRetry() && !request.IsExpired()
}
