// services/dispatch-service/internal/domain/entities/dispatch_request.go
// DispatchRequest entity with matching state machine

package entities

import (
	"fmt"
	"time"
)

// MatchStatus represents match state
type MatchStatus string

const (
	StatusPending      MatchStatus = "pending"
	StatusMatching     MatchStatus = "matching"
	StatusMatched      MatchStatus = "matched"
	StatusAccepted     MatchStatus = "accepted"
	StatusRejected     MatchStatus = "rejected"
	StatusExpired      MatchStatus = "expired"
	StatusCompleted    MatchStatus = "completed"
	StatusFailed       MatchStatus = "failed"
	StatusCancelled    MatchStatus = "cancelled"
)

// DispatchRequest represents a ride matching request
type DispatchRequest struct {
	ID                    string
	RideID                string
	RiderID               string
	PickupLatitude        float64
	PickupLongitude       float64
	DropoffLatitude       float64
	DropoffLongitude      float64
	Status                MatchStatus
	MatchedDriverID       *string // Nil until matched
	ProposedDriverIDs     []string // Top candidates
	RequestedAt           time.Time
	MatchingStartedAt     *time.Time
	MatchedAt             *time.Time
	AcceptedAt            *time.Time
	RejectedAt            *time.Time
	ExpiredAt             *time.Time
	RejectionReason       string
	ExpiryTime            time.Time
	SearchRadius          float64 // kilometers
	MaxSearchRadius       float64 // Maximum allowed
	AttemptCount          int
	MaxAttempts           int
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time

	// Audit fields
	CreatedBy string
	UpdatedBy string
}

// NewDispatchRequest creates a new dispatch request
func NewDispatchRequest(
	rideID, riderID string,
	pickupLat, pickupLng, dropoffLat, dropoffLng float64,
	searchRadius, maxSearchRadius float64,
	maxAttempts int,
) (*DispatchRequest, error) {
	if rideID == "" {
		return nil, fmt.Errorf("ride ID cannot be empty")
	}
	if riderID == "" {
		return nil, fmt.Errorf("rider ID cannot be empty")
	}
	if searchRadius <= 0 {
		return nil, fmt.Errorf("search radius must be positive")
	}

	now := time.Now()
	return &DispatchRequest{
		ID:              fmt.Sprintf("dispatch_%d", now.UnixNano()),
		RideID:          rideID,
		RiderID:         riderID,
		PickupLatitude:  pickupLat,
		PickupLongitude: pickupLng,
		DropoffLatitude: dropoffLat,
		DropoffLongitude: dropoffLng,
		Status:          StatusPending,
		RequestedAt:     now,
		SearchRadius:    searchRadius,
		MaxSearchRadius: maxSearchRadius,
		AttemptCount:    0,
		MaxAttempts:     maxAttempts,
		ExpiryTime:      now.Add(5 * time.Minute),
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// IsValid checks if dispatch request is valid
func (dr *DispatchRequest) IsValid() bool {
	return dr.ID != "" &&
		dr.RideID != "" &&
		dr.RiderID != "" &&
		dr.Status != "" &&
		dr.SearchRadius > 0 &&
		!dr.IsDeleted()
}

// IsDeleted checks if request is deleted
func (dr *DispatchRequest) IsDeleted() bool {
	return dr.DeletedAt != nil
}

// IsExpired checks if request has expired
func (dr *DispatchRequest) IsExpired() bool {
	return time.Now().After(dr.ExpiryTime)
}

// CanStartMatching checks if matching can start
func (dr *DispatchRequest) CanStartMatching() bool {
	return dr.Status == StatusPending && !dr.IsExpired()
}

// StartMatching starts the matching process
func (dr *DispatchRequest) StartMatching() error {
	if !dr.CanStartMatching() {
		return fmt.Errorf("cannot start matching in current state: %s", dr.Status)
	}

	now := time.Now()
	dr.Status = StatusMatching
	dr.MatchingStartedAt = &now
	dr.AttemptCount++
	dr.UpdatedAt = now

	return nil
}

// CanMatch checks if drivers can be matched
func (dr *DispatchRequest) CanMatch() bool {
	return dr.Status == StatusMatching && !dr.IsExpired()
}

// Match sets matched drivers
func (dr *DispatchRequest) Match(driverID string, proposedDrivers []string) error {
	if !dr.CanMatch() {
		return fmt.Errorf("cannot match in current state: %s", dr.Status)
	}

	now := time.Now()
	dr.Status = StatusMatched
	dr.MatchedDriverID = &driverID
	dr.ProposedDriverIDs = proposedDrivers
	dr.MatchedAt = &now
	dr.UpdatedAt = now

	return nil
}

// CanAccept checks if match can be accepted
func (dr *DispatchRequest) CanAccept() bool {
	return dr.Status == StatusMatched && dr.MatchedDriverID != nil
}

// Accept marks match as accepted
func (dr *DispatchRequest) Accept(driverID string) error {
	if !dr.CanAccept() {
		return fmt.Errorf("cannot accept match in current state: %s", dr.Status)
	}
	if dr.MatchedDriverID == nil || *dr.MatchedDriverID != driverID {
		return fmt.Errorf("driver ID mismatch: expected %s, got %s", *dr.MatchedDriverID, driverID)
	}

	now := time.Now()
	dr.Status = StatusAccepted
	dr.AcceptedAt = &now
	dr.UpdatedAt = now

	return nil
}

// CanReject checks if match can be rejected
func (dr *DispatchRequest) CanReject() bool {
	return dr.Status == StatusMatched
}

// Reject rejects the match
func (dr *DispatchRequest) Reject(reason string) error {
	if !dr.CanReject() {
		return fmt.Errorf("cannot reject match in current state: %s", dr.Status)
	}

	now := time.Now()
	dr.Status = StatusRejected
	dr.RejectionReason = reason
	dr.RejectedAt = &now
	dr.MatchedDriverID = nil
	dr.UpdatedAt = now

	return nil
}

// CanRetry checks if matching can be retried
func (dr *DispatchRequest) CanRetry() bool {
	return dr.AttemptCount < dr.MaxAttempts && !dr.IsExpired()
}

// RetryMatching resets for retry
func (dr *DispatchRequest) RetryMatching() error {
	if !dr.CanRetry() {
		if dr.IsExpired() {
			return fmt.Errorf("matching request expired")
		}
		return fmt.Errorf("max retry attempts exceeded: %d/%d", dr.AttemptCount, dr.MaxAttempts)
	}

	now := time.Now()
	dr.Status = StatusPending
	dr.MatchedDriverID = nil
	dr.ProposedDriverIDs = nil
	dr.MatchingStartedAt = nil
	dr.MatchedAt = nil
	dr.RejectedAt = nil
	dr.UpdatedAt = now

	return nil
}

// Expire marks request as expired
func (dr *DispatchRequest) Expire() error {
	if !dr.IsExpired() {
		return fmt.Errorf("request not yet expired")
	}

	now := time.Now()
	dr.Status = StatusExpired
	dr.ExpiredAt = &now
	dr.UpdatedAt = now

	return nil
}

// Complete marks request as completed
func (dr *DispatchRequest) Complete() error {
	if dr.Status != StatusAccepted {
		return fmt.Errorf("can only complete accepted matches")
	}

	now := time.Now()
	dr.Status = StatusCompleted
	dr.UpdatedAt = now

	return nil
}

// Cancel cancels the request
func (dr *DispatchRequest) Cancel() error {
	if dr.Status == StatusCompleted || dr.Status == StatusFailed {
		return fmt.Errorf("cannot cancel request in state: %s", dr.Status)
	}

	now := time.Now()
	dr.Status = StatusCancelled
	dr.UpdatedAt = now

	return nil
}

// ExpandSearchRadius expands the search area
func (dr *DispatchRequest) ExpandSearchRadius(newRadius float64) error {
	if newRadius > dr.MaxSearchRadius {
		newRadius = dr.MaxSearchRadius
	}
	if newRadius <= dr.SearchRadius {
		return fmt.Errorf("new radius must be larger than current")
	}

	dr.SearchRadius = newRadius
	dr.UpdatedAt = time.Now()

	return nil
}
