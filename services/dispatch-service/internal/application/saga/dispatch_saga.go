package saga

import (
	"context"

	ridecontracts "github.com/Abdex1/FamGo-platform/packages/event-bus/contracts/ride"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"
)

// DispatchMatchingPort isolates saga orchestration from concrete use case wiring.
type DispatchMatchingPort interface {
	MatchRide(ctx context.Context, input *usecases.MatchRideInput) (*usecases.MatchResult, error)
	CancelDispatch(ctx context.Context, input *usecases.CancelDispatchInput) error
}

// DispatchSagaHandler connects ride lifecycle events to dispatch matching.
type DispatchSagaHandler struct {
	useCases DispatchMatchingPort
}

func NewDispatchSagaHandler(useCases DispatchMatchingPort) *DispatchSagaHandler {
	return &DispatchSagaHandler{useCases: useCases}
}

// HandleRideCreated consumes ride.created.v1 and starts dispatch matching.
func (h *DispatchSagaHandler) HandleRideCreated(
	ctx context.Context,
	event ridecontracts.RideCreated,
	traceID, correlationID, requestID string,
) (*usecases.MatchResult, error) {
	return h.useCases.MatchRide(ctx, &usecases.MatchRideInput{
		RideID:        event.RideID,
		RiderID:       event.RiderID,
		PickupLat:     event.PickupLat,
		PickupLng:     event.PickupLng,
		DropoffLat:    event.DropoffLat,
		DropoffLng:    event.DropoffLng,
		TraceID:       traceID,
		CorrelationID: correlationID,
		RequestID:     requestID,
	})
}

// CompensationCancelDispatch cancels dispatch when upstream ride saga fails.
func (h *DispatchSagaHandler) CompensationCancelDispatch(
	ctx context.Context,
	dispatchRequestID string,
) error {
	return h.useCases.CancelDispatch(ctx, &usecases.CancelDispatchInput{
		DispatchRequestID: dispatchRequestID,
	})
}
