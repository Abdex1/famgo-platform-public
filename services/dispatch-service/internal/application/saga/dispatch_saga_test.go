package saga

import (
	"context"
	"testing"

	ridecontracts "github.com/Abdex1/FamGo-platform/packages/event-bus/contracts/ride"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"
)

type matchingPortStub struct {
	onMatch func(context.Context, *usecases.MatchRideInput)
}

func (s matchingPortStub) MatchRide(ctx context.Context, input *usecases.MatchRideInput) (*usecases.MatchResult, error) {
	if s.onMatch != nil {
		s.onMatch(ctx, input)
	}
	return &usecases.MatchResult{DispatchRequestID: "dispatch-1"}, nil
}

func (matchingPortStub) CancelDispatch(context.Context, *usecases.CancelDispatchInput) error {
	return nil
}

func TestDispatchSagaHandler_HandleRideCreated(t *testing.T) {
	called := false
	handler := NewDispatchSagaHandler(matchingPortStub{
		onMatch: func(_ context.Context, input *usecases.MatchRideInput) {
			called = true
			if input.RideID != "ride-100" || input.RiderID != "rider-100" {
				t.Fatalf("unexpected ride/rider mapping")
			}
			if input.TraceID != "trace" || input.CorrelationID != "corr" || input.RequestID != "req" {
				t.Fatalf("observability IDs not propagated")
			}
		},
	})

	_, err := handler.HandleRideCreated(context.Background(), ridecontracts.RideCreated{
		RideID:     "ride-100",
		RiderID:    "rider-100",
		PickupLat:  9.03,
		PickupLng:  38.74,
		DropoffLat: 9.05,
		DropoffLng: 38.76,
	}, "trace", "corr", "req")
	if err != nil {
		t.Fatalf("HandleRideCreated() error = %v", err)
	}
	if !called {
		t.Fatal("expected saga handler to invoke MatchRide")
	}
}

func TestDispatchSagaHandler_CompensationCancelDispatch(t *testing.T) {
	cancelled := false
	handler := NewDispatchSagaHandler(cancelPortStub{
		onCancel: func(_ context.Context, input *usecases.CancelDispatchInput) {
			cancelled = true
			if input.DispatchRequestID != "dispatch-1" {
				t.Fatalf("unexpected dispatch id %s", input.DispatchRequestID)
			}
		},
	})

	if err := handler.CompensationCancelDispatch(context.Background(), "dispatch-1"); err != nil {
		t.Fatalf("CompensationCancelDispatch() error = %v", err)
	}
	if !cancelled {
		t.Fatal("expected cancel to be invoked")
	}
}

type cancelPortStub struct {
	onCancel func(context.Context, *usecases.CancelDispatchInput)
}

func (s cancelPortStub) MatchRide(context.Context, *usecases.MatchRideInput) (*usecases.MatchResult, error) {
	return nil, nil
}

func (s cancelPortStub) CancelDispatch(ctx context.Context, input *usecases.CancelDispatchInput) error {
	if s.onCancel != nil {
		s.onCancel(ctx, input)
	}
	return nil
}
