package commands

import "time"

type MatchRideCommand struct {
	RideID            string
	RiderID           string
	PickupLat         float64
	PickupLng         float64
	DropoffLat        float64
	DropoffLng        float64
	SearchRadiusKm    float64
	MaxSearchRadiusKm float64
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type AcceptMatchCommand struct {
	DispatchRequestID string
	DriverID          string
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type RejectMatchCommand struct {
	DispatchRequestID string
	DriverID          string
	Reason            string
	CanRetry          bool
}

type CancelDispatchCommand struct {
	DispatchRequestID string
}

type MatchRideResult struct {
	DispatchRequestID string
	MatchedDriverID   string
	ProposedDrivers   []string
	Status            string
	MatchedAt         time.Time
}
