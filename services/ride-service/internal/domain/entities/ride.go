// services/ride-service/internal/domain/entities/ride.go
// Ride domain entity with state machine

package entities

import (
	"fmt"
	"time"
)

// RideStatus represents ride state in state machine
type RideStatus string

const (
	StatusRequested   RideStatus = "requested"
	StatusAccepted    RideStatus = "accepted"
	StatusPickedUp    RideStatus = "picked_up"
	StatusOngoing     RideStatus = "ongoing"
	StatusCompleted   RideStatus = "completed"
	StatusCancelled   RideStatus = "cancelled"
	StatusNoShow      RideStatus = "no_show"
)

// Ride represents a ride request and lifecycle
type Ride struct {
	ID                   string
	RiderID              string
	DriverID             *string // Nil until accepted
	PickupLatitude       float64
	PickupLongitude      float64
	DropoffLatitude      float64
	DropoffLongitude     float64
	EstimatedDistance    float64 // kilometers
	EstimatedDuration    time.Duration
	Status               RideStatus
	RequestedAt          time.Time
	AcceptedAt           *time.Time
	PickedUpAt           *time.Time
	StartedAt            *time.Time
	CompletedAt          *time.Time
	CancelledAt          *time.Time
	CancelReason         string
	RideType             string // "standard", "pool", "xl"
	PassengerCount       int
	EstimatedFare        float64
	ActualFare           *float64
	PaymentMethod        string
	PaymentStatus        string // "pending", "completed", "failed"
	RiderRating          *int
	DriverRating         *int
	PickupAddress        string
	DropoffAddress       string
	SpecialRequests      string
	RoutePolyline        string
	ActualDistance       *float64 // Set when completed
	ActualDuration       *time.Duration
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

// NewRide creates a new ride entity
func NewRide(riderID string, pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*Ride, error) {
	if riderID == "" {
		return nil, fmt.Errorf("rider ID cannot be empty")
	}

	now := time.Now()
	return &Ride{
		ID:                fmt.Sprintf("ride_%d", now.UnixNano()),
		RiderID:           riderID,
		PickupLatitude:    pickupLat,
		PickupLongitude:   pickupLng,
		DropoffLatitude:   dropoffLat,
		DropoffLongitude:  dropoffLng,
		Status:            StatusRequested,
		RequestedAt:       now,
		PassengerCount:    1,
		RideType:          "standard",
		PaymentStatus:     "pending",
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// IsValid checks if ride is valid
func (r *Ride) IsValid() bool {
	return r.ID != "" &&
		r.RiderID != "" &&
		r.Status != "" &&
		r.PickupLatitude >= -90 && r.PickupLatitude <= 90 &&
		r.PickupLongitude >= -180 && r.PickupLongitude <= 180 &&
		r.DropoffLatitude >= -90 && r.DropoffLatitude <= 90 &&
		r.DropoffLongitude >= -180 && r.DropoffLongitude <= 180
}

// CanAccept checks if ride can be accepted
func (r *Ride) CanAccept() bool {
	return r.Status == StatusRequested && r.DeletedAt == nil
}

// Accept marks ride as accepted by driver
func (r *Ride) Accept(driverID string) error {
	if !r.CanAccept() {
		return fmt.Errorf("ride cannot be accepted in current state: %s", r.Status)
	}
	if driverID == "" {
		return fmt.Errorf("driver ID cannot be empty")
	}

	now := time.Now()
	r.DriverID = &driverID
	r.Status = StatusAccepted
	r.AcceptedAt = &now
	r.UpdatedAt = now

	return nil
}

// CanPickup checks if ride can be picked up
func (r *Ride) CanPickup() bool {
	return r.Status == StatusAccepted && r.DriverID != nil && r.DeletedAt == nil
}

// PickupPassenger marks passenger as picked up
func (r *Ride) PickupPassenger() error {
	if !r.CanPickup() {
		return fmt.Errorf("ride cannot be picked up in current state: %s", r.Status)
	}

	now := time.Now()
	r.Status = StatusPickedUp
	r.PickedUpAt = &now
	r.UpdatedAt = now

	return nil
}

// CanStart checks if ride can be started
func (r *Ride) CanStart() bool {
	return r.Status == StatusPickedUp && r.PickedUpAt != nil && r.DeletedAt == nil
}

// StartRide marks ride as started (en route)
func (r *Ride) StartRide() error {
	if !r.CanStart() {
		return fmt.Errorf("ride cannot be started in current state: %s", r.Status)
	}

	now := time.Now()
	r.Status = StatusOngoing
	r.StartedAt = &now
	r.UpdatedAt = now

	return nil
}

// CanComplete checks if ride can be completed
func (r *Ride) CanComplete() bool {
	return r.Status == StatusOngoing && r.StartedAt != nil && r.DeletedAt == nil
}

// CompleteRide marks ride as completed
func (r *Ride) CompleteRide(actualDistance, actualFare float64) error {
	if !r.CanComplete() {
		return fmt.Errorf("ride cannot be completed in current state: %s", r.Status)
	}

	now := time.Now()
	r.Status = StatusCompleted
	r.CompletedAt = &now
	r.ActualDistance = &actualDistance
	r.ActualFare = &actualFare
	r.PaymentStatus = "completed"

	duration := now.Sub(*r.StartedAt)
	r.ActualDuration = &duration

	r.UpdatedAt = now

	return nil
}

// CanCancel checks if ride can be cancelled
func (r *Ride) CanCancel() bool {
	return (r.Status == StatusRequested ||
		r.Status == StatusAccepted ||
		r.Status == StatusPickedUp) &&
		r.DeletedAt == nil
}

// Cancel cancels the ride
func (r *Ride) Cancel(reason string) error {
	if !r.CanCancel() {
		return fmt.Errorf("ride cannot be cancelled in current state: %s", r.Status)
	}

	now := time.Now()
	r.Status = StatusCancelled
	r.CancelledAt = &now
	r.CancelReason = reason
	r.UpdatedAt = now

	return nil
}

// SetNoShow marks ride as no-show
func (r *Ride) SetNoShow() error {
	if r.Status != StatusAccepted && r.Status != StatusPickedUp {
		return fmt.Errorf("ride cannot be no-show in current state: %s", r.Status)
	}

	now := time.Now()
	r.Status = StatusNoShow
	r.CancelledAt = &now
	r.CancelReason = "passenger_no_show"
	r.UpdatedAt = now

	return nil
}

// SetRiderRating sets rider rating
func (r *Ride) SetRiderRating(rating int) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	if r.Status != StatusCompleted {
		return fmt.Errorf("can only rate completed rides")
	}

	r.RiderRating = &rating
	r.UpdatedAt = time.Now()

	return nil
}

// SetDriverRating sets driver rating
func (r *Ride) SetDriverRating(rating int) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	if r.Status != StatusCompleted {
		return fmt.Errorf("can only rate completed rides")
	}

	r.DriverRating = &rating
	r.UpdatedAt = time.Now()

	return nil
}

// IsCompleted checks if ride is completed
func (r *Ride) IsCompleted() bool {
	return r.Status == StatusCompleted || r.Status == StatusCancelled || r.Status == StatusNoShow
}

// IsCancelled checks if ride is cancelled
func (r *Ride) IsCancelled() bool {
	return r.Status == StatusCancelled || r.Status == StatusNoShow
}

// GetDuration returns ride duration in minutes
func (r *Ride) GetDuration() float64 {
	if r.ActualDuration != nil {
		return r.ActualDuration.Minutes()
	}
	if r.EstimatedDuration > 0 {
		return r.EstimatedDuration.Minutes()
	}
	return 0
}

// GetDistance returns ride distance in kilometers
func (r *Ride) GetDistance() float64 {
	if r.ActualDistance != nil {
		return *r.ActualDistance
	}
	return r.EstimatedDistance
}

// UpdateEstimate updates estimated distance and duration
func (r *Ride) UpdateEstimate(distance float64, duration time.Duration) {
	r.EstimatedDistance = distance
	r.EstimatedDuration = duration
	r.UpdatedAt = time.Now()
}
