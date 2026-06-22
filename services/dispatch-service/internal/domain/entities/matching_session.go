package entities

import (
	"fmt"
	"time"
)

// MatchingSessionStatus tracks an active matching attempt.
type MatchingSessionStatus string

const (
	MatchingSessionActive    MatchingSessionStatus = "active"
	MatchingSessionMatched   MatchingSessionStatus = "matched"
	MatchingSessionFailed    MatchingSessionStatus = "failed"
	MatchingSessionExpired   MatchingSessionStatus = "expired"
	MatchingSessionCancelled MatchingSessionStatus = "cancelled"
)

// MatchingSession records a single matching attempt for audit and saga correlation.
type MatchingSession struct {
	ID                string
	DispatchRequestID string
	RideID            string
	Status            MatchingSessionStatus
	Algorithm         string
	SearchRadiusKm    float64
	CandidateCount    int
	SelectedDriverID  *string
	StartedAt         time.Time
	CompletedAt       *time.Time
	ExpiresAt         time.Time
	FailureReason     string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewMatchingSession(
	dispatchRequestID, rideID, algorithm string,
	searchRadiusKm float64,
	timeout time.Duration,
) (*MatchingSession, error) {
	if dispatchRequestID == "" || rideID == "" {
		return nil, fmt.Errorf("dispatch request ID and ride ID are required")
	}
	if algorithm == "" {
		algorithm = "nearest_available_v1"
	}
	now := time.Now().UTC()
	return &MatchingSession{
		ID:                fmt.Sprintf("match_session_%d", now.UnixNano()),
		DispatchRequestID: dispatchRequestID,
		RideID:            rideID,
		Status:            MatchingSessionActive,
		Algorithm:         algorithm,
		SearchRadiusKm:    searchRadiusKm,
		StartedAt:         now,
		ExpiresAt:         now.Add(timeout),
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (s *MatchingSession) Complete(driverID string) error {
	if s.Status != MatchingSessionActive {
		return fmt.Errorf("cannot complete session in status %s", s.Status)
	}
	now := time.Now().UTC()
	s.Status = MatchingSessionMatched
	s.SelectedDriverID = &driverID
	s.CompletedAt = &now
	s.UpdatedAt = now
	return nil
}

func (s *MatchingSession) Fail(reason string) error {
	now := time.Now().UTC()
	s.Status = MatchingSessionFailed
	s.FailureReason = reason
	s.CompletedAt = &now
	s.UpdatedAt = now
	return nil
}

func (s *MatchingSession) MarkExpired() {
	now := time.Now().UTC()
	s.Status = MatchingSessionExpired
	s.CompletedAt = &now
	s.UpdatedAt = now
}

func (s *MatchingSession) Cancel() {
	now := time.Now().UTC()
	s.Status = MatchingSessionCancelled
	s.CompletedAt = &now
	s.UpdatedAt = now
}

func (s *MatchingSession) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

func (s *MatchingSession) SetCandidateCount(count int) {
	s.CandidateCount = count
	s.UpdatedAt = time.Now().UTC()
}
