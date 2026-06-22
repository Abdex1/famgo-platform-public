package ports

import "context"

// DispatchEventPublisher publishes dispatch domain integration events.
type DispatchEventPublisher interface {
	PublishMatchingStarted(ctx context.Context, event MatchingStartedEvent) error
	PublishDriverMatched(ctx context.Context, event DriverMatchedEvent) error
	PublishDriverAssigned(ctx context.Context, event DriverAssignedEvent) error
	PublishMatchingFailed(ctx context.Context, event MatchingFailedEvent) error
	PublishMatchingExpired(ctx context.Context, event MatchingExpiredEvent) error
}

type MatchingStartedEvent struct {
	DispatchRequestID string
	RideID            string
	RiderID           string
	SearchRadiusKm    float64
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type DriverMatchedEvent struct {
	DispatchRequestID string
	RideID            string
	DriverID          string
	ProposedDrivers   []string
	MatchScore        float64
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type DriverAssignedEvent struct {
	DispatchRequestID string
	RideID            string
	DriverID          string
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type MatchingFailedEvent struct {
	DispatchRequestID string
	RideID            string
	Reason            string
	AttemptCount      int
	TraceID           string
	CorrelationID     string
	RequestID         string
}

type MatchingExpiredEvent struct {
	DispatchRequestID string
	RideID            string
	TraceID           string
	CorrelationID     string
	RequestID         string
}
