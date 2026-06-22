// services/ride-service/internal/domain/entities/ride_test.go
package entities

import (
	"testing"
	"time"
)

func TestNewRide(t *testing.T) {
	ride, err := NewRide("rider1", 9.0320, 38.7469, 9.0350, 38.7500)
	
	if err != nil {
		t.Errorf("NewRide() error = %v", err)
	}
	
	if ride == nil {
		t.Fatalf("NewRide() returned nil")
	}
	
	if ride.RiderID != "rider1" {
		t.Errorf("RiderID = %s, want rider1", ride.RiderID)
	}
	
	if ride.Status != StatusRequested {
		t.Errorf("Status = %v, want %v", ride.Status, StatusRequested)
	}
}

func TestRideStateMachine(t *testing.T) {
	ride, _ := NewRide("rider1", 9.0320, 38.7469, 9.0350, 38.7500)
	
	// Test valid transitions
	if !ride.CanAccept() {
		t.Errorf("Ride should be acceptable in REQUESTED state")
	}
	
	if err := ride.Accept("driver1"); err != nil {
		t.Errorf("Accept() error = %v", err)
	}
	
	if ride.Status != StatusAccepted {
		t.Errorf("Status = %v, want %v", ride.Status, StatusAccepted)
	}
	
	// Pickup
	if !ride.CanPickup() {
		t.Errorf("Ride should be pickupable in ACCEPTED state")
	}
	
	if err := ride.PickupPassenger(); err != nil {
		t.Errorf("PickupPassenger() error = %v", err)
	}
	
	if ride.Status != StatusPickedUp {
		t.Errorf("Status = %v, want %v", ride.Status, StatusPickedUp)
	}
	
	// Start
	if !ride.CanStart() {
		t.Errorf("Ride should be startable in PICKED_UP state")
	}
	
	if err := ride.StartRide(); err != nil {
		t.Errorf("StartRide() error = %v", err)
	}
	
	if ride.Status != StatusOngoing {
		t.Errorf("Status = %v, want %v", ride.Status, StatusOngoing)
	}
	
	// Complete
	if !ride.CanComplete() {
		t.Errorf("Ride should be completable in ONGOING state")
	}
	
	if err := ride.CompleteRide(5.0, 150.0); err != nil {
		t.Errorf("CompleteRide() error = %v", err)
	}
	
	if ride.Status != StatusCompleted {
		t.Errorf("Status = %v, want %v", ride.Status, StatusCompleted)
	}
}

func TestRideCancel(t *testing.T) {
	ride, _ := NewRide("rider1", 9.0320, 38.7469, 9.0350, 38.7500)
	
	if err := ride.Accept("driver1"); err != nil {
		t.Fatalf("Accept() error = %v", err)
	}
	
	if !ride.CanCancel() {
		t.Errorf("Ride should be cancellable in ACCEPTED state")
	}
	
	if err := ride.Cancel("rider_cancelled"); err != nil {
		t.Errorf("Cancel() error = %v", err)
	}
	
	if ride.Status != StatusCancelled {
		t.Errorf("Status = %v, want %v", ride.Status, StatusCancelled)
	}
	
	if ride.CancelReason != "rider_cancelled" {
		t.Errorf("CancelReason = %s, want rider_cancelled", ride.CancelReason)
	}
}

func TestRideRating(t *testing.T) {
	ride, _ := NewRide("rider1", 9.0320, 38.7469, 9.0350, 38.7500)
	ride.Status = StatusCompleted
	
	if err := ride.SetDriverRating(5); err != nil {
		t.Errorf("SetDriverRating() error = %v", err)
	}
	
	if ride.DriverRating == nil || *ride.DriverRating != 5 {
		t.Errorf("DriverRating = %v, want 5", ride.DriverRating)
	}
	
	// Invalid rating
	if err := ride.SetDriverRating(6); err == nil {
		t.Errorf("SetDriverRating(6) should error")
	}
}

func TestRideMetrics(t *testing.T) {
	ride, _ := NewRide("rider1", 9.0320, 38.7469, 9.0350, 38.7500)
	ride.EstimatedDistance = 10.0
	ride.EstimatedDuration = 30 * time.Minute
	ride.EstimatedFare = 150.0
	
	if distance := ride.GetDistance(); distance != 10.0 {
		t.Errorf("GetDistance() = %f, want 10.0", distance)
	}
	
	if duration := ride.GetDuration(); duration != 30.0 {
		t.Errorf("GetDuration() = %f, want 30.0", duration)
	}
}
