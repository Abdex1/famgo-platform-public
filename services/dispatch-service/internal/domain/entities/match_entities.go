package entities

import (
	"time"
)

// MatchRequest represents a dispatch matching request
type MatchRequest struct {
	ID                string
	RideID            string
	PickupLat         float64
	PickupLng         float64
	DropoffLat        float64
	DropoffLng        float64
	RiderRating       float64
	RequestedVehicleType string
	Status            string // PENDING, MATCHED, FAILED, EXPIRED
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExpiresAt         time.Time
}

func NewMatchRequest(
	id string,
	rideID string,
	pickupLat, pickupLng float64,
	dropoffLat, dropoffLng float64,
	vehicleType string,
	timeoutSeconds int32,
) *MatchRequest {
	now := time.Now().UTC()
	return &MatchRequest{
		ID:                   id,
		RideID:               rideID,
		PickupLat:            pickupLat,
		PickupLng:            pickupLng,
		DropoffLat:           dropoffLat,
		DropoffLng:           dropoffLng,
		RequestedVehicleType: vehicleType,
		Status:               "PENDING",
		CreatedAt:            now,
		UpdatedAt:            now,
		ExpiresAt:            now.Add(time.Duration(timeoutSeconds) * time.Second),
	}
}

// IsExpired checks if match request has timed out
func (mr *MatchRequest) IsExpired() bool {
	return time.Now().UTC().After(mr.ExpiresAt)
}

// DriverMatch represents a driver candidate for matching
type DriverMatch struct {
	DriverID        string
	Latitude        float64
	Longitude       float64
	Distance        float64       // km
	Rating          float64       // 1-5 stars
	AcceptedRides   int32
	CancelledRides  int32
	OnlineStatus    bool
	Score           float64       // Composite matching score (0-100)
}

// MatchResult represents the result of a matching operation
type MatchResult struct {
	MatchRequestID string
	RideID         string
	SelectedDriverID string
	Score          float64
	ETA            int32         // minutes
	Distance       float64       // km
	Confidence     float64       // 0-1 confidence level
	CreatedAt      time.Time
}

func NewMatchResult(
	matchReqID string,
	rideID string,
	driverID string,
	score float64,
	eta int32,
	distance float64,
	confidence float64,
) *MatchResult {
	return &MatchResult{
		MatchRequestID:   matchReqID,
		RideID:           rideID,
		SelectedDriverID: driverID,
		Score:            score,
		ETA:              eta,
		Distance:         distance,
		Confidence:       confidence,
		CreatedAt:        time.Now().UTC(),
	}
}
