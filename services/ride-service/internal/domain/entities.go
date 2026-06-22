// services/ride-service/internal/domain/entities.go
// Ride Service Domain Entities with State Machine

package domain

import (
	"time"
)

// RideStatus represents ride lifecycle state
type RideStatus string

const (
	RideStatusRequested      RideStatus = "REQUESTED"       // User creates ride
	RideStatusSearching      RideStatus = "SEARCHING"       // Dispatch searching
	RideStatusAssigned       RideStatus = "ASSIGNED"        // Driver assigned
	RideStatusDriverArriving RideStatus = "DRIVER_ARRIVING" // Driver en route
	RideStatusStarted        RideStatus = "STARTED"         // Pickup complete
	RideStatusCompleted      RideStatus = "COMPLETED"       // Dropoff complete
	RideStatusCancelled      RideStatus = "CANCELLED"       // Cancelled
)

// Ride represents a ride request and lifecycle
type Ride struct {
	ID              string    // UUID
	PassengerID     string    // Foreign key to user
	DriverID        string    // Assigned driver (nullable)
	PickupLat       float64   // Pickup latitude
	PickupLon       float64   // Pickup longitude
	DropoffLat      float64   // Dropoff latitude
	DropoffLon      float64   // Dropoff longitude
	Status          RideStatus
	EstimatedFare   float32   // Calculated by pricing service
	ActualFare      float32   // Final fare after completion
	PickupTime      *time.Time // When pickup started
	DropoffTime     *time.Time // When dropoff completed
	CancellationReason string  // If cancelled
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// RideStatusHistory tracks status transitions
type RideStatusHistory struct {
	ID        string    // UUID
	RideID    string    // Foreign key
	OldStatus RideStatus
	NewStatus RideStatus
	ChangedAt time.Time
}

// NewRideWithID creates a new ride in REQUESTED state
// ID is provided by application/factory layer to avoid external dependencies
func NewRideWithID(id, passengerID string, pickupLat, pickupLon, dropoffLat, dropoffLon float64) *Ride {
	return &Ride{
		ID:          id,
		PassengerID: passengerID,
		PickupLat:   pickupLat,
		PickupLon:   pickupLon,
		DropoffLat:  dropoffLat,
		DropoffLon:  dropoffLon,
		Status:      RideStatusRequested,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CanTransitionTo validates state transitions
func (r *Ride) CanTransitionTo(newStatus RideStatus) bool {
	// Define allowed state transitions
	allowedTransitions := map[RideStatus][]RideStatus{
		RideStatusRequested: {
			RideStatusSearching,
			RideStatusCancelled,
		},
		RideStatusSearching: {
			RideStatusAssigned,
			RideStatusCancelled,
		},
		RideStatusAssigned: {
			RideStatusDriverArriving,
			RideStatusCancelled,
		},
		RideStatusDriverArriving: {
			RideStatusStarted,
			RideStatusCancelled,
		},
		RideStatusStarted: {
			RideStatusCompleted,
			RideStatusCancelled,
		},
		RideStatusCompleted:  {},  // Terminal state
		RideStatusCancelled:  {},  // Terminal state
	}

	for _, allowed := range allowedTransitions[r.Status] {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// TransitionTo performs state transition with validation
func (r *Ride) TransitionTo(newStatus RideStatus) error {
	if !r.CanTransitionTo(newStatus) {
		return ErrInvalidStateTransition
	}
	r.Status = newStatus
	r.UpdatedAt = time.Now()
	return nil
}

// AssignDriver assigns a driver to ride
func (r *Ride) AssignDriver(driverID string) {
	r.DriverID = driverID
	r.UpdatedAt = time.Now()
}

// SetEstimatedFare sets estimated fare from pricing service
func (r *Ride) SetEstimatedFare(fare float32) {
	r.EstimatedFare = fare
	r.UpdatedAt = time.Now()
}

// StartPickup marks when driver starts pickup
func (r *Ride) StartPickup() {
	now := time.Now()
	r.PickupTime = &now
	r.UpdatedAt = time.Now()
}

// CompleteRide marks ride as completed with final fare
func (r *Ride) CompleteRide(actualFare float32) {
	now := time.Now()
	r.DropoffTime = &now
	r.ActualFare = actualFare
	r.UpdatedAt = time.Now()
}

// CancelRide marks ride as cancelled with reason
func (r *Ride) CancelRide(reason string) {
	r.CancellationReason = reason
	r.UpdatedAt = time.Now()
}

// IsTerminalState checks if ride is in terminal state
func (r *Ride) IsTerminalState() bool {
	return r.Status == RideStatusCompleted || r.Status == RideStatusCancelled
}

// IsActive checks if ride is active (not completed/cancelled)
func (r *Ride) IsActive() bool {
	return !r.IsTerminalState()
}
