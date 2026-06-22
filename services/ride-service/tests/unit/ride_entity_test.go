// services/ride-service/tests/unit/ride_entity_test.go
// Unit tests for Ride entity

package unit

import (
	"testing"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

func TestNewRide(t *testing.T) {
	passengerID := "passenger123"
	pickupLat, pickupLon := 37.7749, -122.4194
	dropoffLat, dropoffLon := 37.8044, -122.2712

	ride := domain.NewRide(passengerID, pickupLat, pickupLon, dropoffLat, dropoffLon)

	if ride.ID == "" {
		t.Error("ride ID should not be empty")
	}
	if ride.PassengerID != passengerID {
		t.Errorf("expected passenger ID %s, got %s", passengerID, ride.PassengerID)
	}
	if ride.Status != domain.RideStatusRequested {
		t.Errorf("expected status %s, got %s", domain.RideStatusRequested, ride.Status)
	}
	if ride.PickupLat != pickupLat {
		t.Errorf("expected pickup lat %f, got %f", pickupLat, ride.PickupLat)
	}
}

func TestRideStateTransitions(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)

	tests := []struct {
		name          string
		currentStatus domain.RideStatus
		newStatus     domain.RideStatus
		shouldAllow   bool
	}{
		{"Requested -> Searching", domain.RideStatusRequested, domain.RideStatusSearching, true},
		{"Requested -> Cancelled", domain.RideStatusRequested, domain.RideStatusCancelled, true},
		{"Requested -> Completed", domain.RideStatusRequested, domain.RideStatusCompleted, false},
		{"Searching -> Assigned", domain.RideStatusSearching, domain.RideStatusAssigned, true},
		{"Assigned -> DriverArriving", domain.RideStatusAssigned, domain.RideStatusDriverArriving, true},
		{"DriverArriving -> Started", domain.RideStatusDriverArriving, domain.RideStatusStarted, true},
		{"Started -> Completed", domain.RideStatusStarted, domain.RideStatusCompleted, true},
		{"Completed -> Cancelled", domain.RideStatusCompleted, domain.RideStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ride.Status = tt.currentStatus
			allowed := ride.CanTransitionTo(tt.newStatus)
			if allowed != tt.shouldAllow {
				t.Errorf("expected %v, got %v", tt.shouldAllow, allowed)
			}
		})
	}
}

func TestAssignDriver(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
	driverID := "driver456"

	ride.AssignDriver(driverID)

	if ride.DriverID != driverID {
		t.Errorf("expected driver ID %s, got %s", driverID, ride.DriverID)
	}
}

func TestSetEstimatedFare(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
	fare := float32(15.50)

	ride.SetEstimatedFare(fare)

	if ride.EstimatedFare != fare {
		t.Errorf("expected fare %f, got %f", fare, ride.EstimatedFare)
	}
}

func TestStartPickup(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)

	beforeTime := time.Now()
	ride.StartPickup()
	afterTime := time.Now()

	if ride.PickupTime == nil {
		t.Error("pickup time should be set")
	}
	if ride.PickupTime.Before(beforeTime) || ride.PickupTime.After(afterTime) {
		t.Error("pickup time should be close to now")
	}
}

func TestCompleteRide(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
	actualFare := float32(18.75)

	beforeTime := time.Now()
	ride.CompleteRide(actualFare)
	afterTime := time.Now()

	if ride.DropoffTime == nil {
		t.Error("dropoff time should be set")
	}
	if ride.ActualFare != actualFare {
		t.Errorf("expected actual fare %f, got %f", actualFare, ride.ActualFare)
	}
	if ride.DropoffTime.Before(beforeTime) || ride.DropoffTime.After(afterTime) {
		t.Error("dropoff time should be close to now")
	}
}

func TestCancelRide(t *testing.T) {
	ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
	reason := "passenger request"

	ride.CancelRide(reason)

	if ride.CancellationReason != reason {
		t.Errorf("expected reason %s, got %s", reason, ride.CancellationReason)
	}
}

func TestIsTerminalState(t *testing.T) {
	tests := []struct {
		status      domain.RideStatus
		isTerminal  bool
	}{
		{domain.RideStatusRequested, false},
		{domain.RideStatusSearching, false},
		{domain.RideStatusAssigned, false},
		{domain.RideStatusStarted, false},
		{domain.RideStatusCompleted, true},
		{domain.RideStatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
			ride.Status = tt.status
			if ride.IsTerminalState() != tt.isTerminal {
				t.Errorf("expected terminal state %v, got %v", tt.isTerminal, ride.IsTerminalState())
			}
		})
	}
}

func TestIsActive(t *testing.T) {
	tests := []struct {
		status    domain.RideStatus
		isActive  bool
	}{
		{domain.RideStatusRequested, true},
		{domain.RideStatusSearching, true},
		{domain.RideStatusCompleted, false},
		{domain.RideStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			ride := domain.NewRide("passenger123", 37.7749, -122.4194, 37.8044, -122.2712)
			ride.Status = tt.status
			if ride.IsActive() != tt.isActive {
				t.Errorf("expected active %v, got %v", tt.isActive, ride.IsActive())
			}
		})
	}
}
