// services/gps-service/internal/domain/entities/driver_location.go
// Driver location entity with status tracking and history

package entities

import (
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// DriverStatus represents driver's availability status
type DriverStatus string

const (
	StatusOnline      DriverStatus = "online"
	StatusOffline     DriverStatus = "offline"
	StatusOnRide      DriverStatus = "on_ride"
	StatusBreak       DriverStatus = "break"
	StatusMaintenance DriverStatus = "maintenance"
)

// DriverLocation is the core domain entity for driver locations
type DriverLocation struct {
	ID                    string
	DriverID              string
	CurrentLocation       *valueobjects.Geolocation
	PreviousLocation      *valueobjects.Geolocation
	Status                DriverStatus
	IsOnline              bool
	LastUpdateAt          time.Time
	LastSyncAt            time.Time
	LocationAccuracy      float64 // meters
	Speed                 float64 // m/s
	Heading               float64 // degrees
	LastSeenAt            time.Time
	ConsecutiveFailures   int
	VehicleID             string
	VehicleRegistration   string
	ServiceStatus         string // "active", "inactive", "cancelled"
	AcceptedRideCount     int
	CompletedRideCount    int
	CancelledRideCount    int
	AverageRating         float64
	RecentAcceptanceRate  float64 // percentage
	GeohashPrefix         string  // For indexing (e.g., "wx123")
	LastLocationHash      string  // For detecting duplicates
	IsVerified            bool
	IsDocumentsExpired    bool
	IsBanned              bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time

	// Audit fields
	CreatedBy string
	UpdatedBy string
}

// NewDriverLocation creates a new driver location entity
func NewDriverLocation(driverID string, location *valueobjects.Geolocation) (*DriverLocation, error) {
	if driverID == "" {
		return nil, fmt.Errorf("driver ID cannot be empty")
	}
	if location == nil || !location.IsValid() {
		return nil, fmt.Errorf("location must be valid")
	}

	now := time.Now()
	return &DriverLocation{
		ID:               fmt.Sprintf("%s-%d", driverID, now.UnixNano()),
		DriverID:         driverID,
		CurrentLocation:  location,
		Status:           StatusOnline,
		IsOnline:         true,
		LastUpdateAt:     now,
		LastSyncAt:       now,
		LocationAccuracy: location.Accuracy,
		Speed:            location.Speed,
		Heading:          location.Heading,
		LastSeenAt:       now,
		ConsecutiveFailures: 0,
		ServiceStatus:    "active",
		AverageRating:    5.0,
		RecentAcceptanceRate: 100.0,
		IsVerified:       true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsValid checks if the driver location entity is valid
func (dl *DriverLocation) IsValid() bool {
	return dl.ID != "" &&
		dl.DriverID != "" &&
		dl.CurrentLocation != nil &&
		dl.CurrentLocation.IsValid() &&
		dl.Status != "" &&
		!dl.IsDeleted()
}

// IsActive checks if driver is actively available
func (dl *DriverLocation) IsActive() bool {
	return dl.IsOnline &&
		dl.Status == StatusOnline &&
		!dl.IsBanned &&
		!dl.IsDocumentsExpired &&
		dl.IsVerified
}

// IsDeleted checks if driver is deleted
func (dl *DriverLocation) IsDeleted() bool {
	return dl.DeletedAt != nil
}

// DistanceTo calculates distance to another driver's location in kilometers
func (dl *DriverLocation) DistanceTo(other *DriverLocation) float64 {
	if other == nil || other.CurrentLocation == nil {
		return 0
	}
	return dl.CurrentLocation.DistanceToKm(other.CurrentLocation)
}

// BearingTo calculates bearing to another driver
func (dl *DriverLocation) BearingTo(other *DriverLocation) float64 {
	if other == nil || other.CurrentLocation == nil {
		return 0
	}
	return dl.CurrentLocation.BearingTo(other.CurrentLocation)
}

// ETATo calculates estimated arrival time to another driver in minutes
func (dl *DriverLocation) ETATo(other *DriverLocation, baseSpeedKmH float64) float64 {
	if other == nil || other.CurrentLocation == nil {
		return 0
	}
	return dl.CurrentLocation.EstimatedArrivalTime(other.CurrentLocation, baseSpeedKmH)
}

// IsWithinRadius checks if another driver is within specified radius in km
func (dl *DriverLocation) IsWithinRadius(other *DriverLocation, radiusKm float64) bool {
	if other == nil || other.CurrentLocation == nil {
		return false
	}
	return dl.CurrentLocation.IsWithinRadius(other.CurrentLocation, radiusKm)
}

// UpdateLocation updates current location and shifts previous location
func (dl *DriverLocation) UpdateLocation(newLocation *valueobjects.Geolocation) error {
	if newLocation == nil || !newLocation.IsValid() {
		return fmt.Errorf("new location must be valid")
	}

	dl.PreviousLocation = dl.CurrentLocation
	dl.CurrentLocation = newLocation
	dl.LastUpdateAt = time.Now()
	dl.UpdatedAt = time.Now()
	dl.LocationAccuracy = newLocation.Accuracy
	dl.Speed = newLocation.Speed
	dl.Heading = newLocation.Heading
	dl.ConsecutiveFailures = 0

	return nil
}

// GoOnline marks driver as online
func (dl *DriverLocation) GoOnline() {
	dl.IsOnline = true
	dl.Status = StatusOnline
	dl.LastSeenAt = time.Now()
	dl.UpdatedAt = time.Now()
	dl.ConsecutiveFailures = 0
}

// GoOffline marks driver as offline
func (dl *DriverLocation) GoOffline() {
	dl.IsOnline = false
	dl.Status = StatusOffline
	dl.UpdatedAt = time.Now()
}

// GoOnRide marks driver as on a ride
func (dl *DriverLocation) GoOnRide() {
	dl.IsOnline = true
	dl.Status = StatusOnRide
	dl.UpdatedAt = time.Now()
}

// GoOnBreak marks driver as on break
func (dl *DriverLocation) GoOnBreak() {
	dl.IsOnline = false
	dl.Status = StatusBreak
	dl.UpdatedAt = time.Now()
}

// GoToMaintenance marks driver as in maintenance
func (dl *DriverLocation) GoToMaintenance() {
	dl.IsOnline = false
	dl.Status = StatusMaintenance
	dl.UpdatedAt = time.Now()
}

// SetStatus sets driver status
func (dl *DriverLocation) SetStatus(status DriverStatus) error {
	if status != StatusOnline &&
		status != StatusOffline &&
		status != StatusOnRide &&
		status != StatusBreak &&
		status != StatusMaintenance {
		return fmt.Errorf("invalid status: %s", status)
	}
	dl.Status = status
	dl.UpdatedAt = time.Now()
	return nil
}

// RecordFailure records a location update failure
func (dl *DriverLocation) RecordFailure() {
	dl.ConsecutiveFailures++
	dl.UpdatedAt = time.Now()
}

// ResetFailures resets consecutive failures counter
func (dl *DriverLocation) ResetFailures() {
	dl.ConsecutiveFailures = 0
}

// UpdateRideStats updates ride-related statistics
func (dl *DriverLocation) UpdateRideStats(accepted, completed, cancelled int) {
	dl.AcceptedRideCount += accepted
	dl.CompletedRideCount += completed
	dl.CancelledRideCount += cancelled

	// Calculate acceptance rate
	total := dl.AcceptedRideCount + dl.CancelledRideCount
	if total > 0 {
		dl.RecentAcceptanceRate = (float64(dl.AcceptedRideCount) / float64(total)) * 100
	}

	dl.UpdatedAt = time.Now()
}

// UpdateRating updates average driver rating
func (dl *DriverLocation) UpdateRating(newRating float64, totalRatings int) {
	if totalRatings > 0 {
		dl.AverageRating = (dl.AverageRating*float64(totalRatings-1) + newRating) / float64(totalRatings)
	}
	dl.UpdatedAt = time.Now()
}

// Ban bans the driver
func (dl *DriverLocation) Ban() {
	dl.IsBanned = true
	dl.IsOnline = false
	dl.Status = StatusOffline
	dl.UpdatedAt = time.Now()
}

// Unban unbans the driver
func (dl *DriverLocation) Unban() {
	dl.IsBanned = false
	dl.UpdatedAt = time.Now()
}

// MarkDocumentsExpired marks driver documents as expired
func (dl *DriverLocation) MarkDocumentsExpired() {
	dl.IsDocumentsExpired = true
	dl.IsOnline = false
	dl.Status = StatusOffline
	dl.UpdatedAt = time.Now()
}

// RenewDocuments marks driver documents as renewed
func (dl *DriverLocation) RenewDocuments() {
	dl.IsDocumentsExpired = false
	dl.UpdatedAt = time.Now()
}

// IsStale checks if driver location is stale (no updates for specified duration)
func (dl *DriverLocation) IsStale(staleDuration time.Duration) bool {
	return time.Since(dl.LastUpdateAt) > staleDuration
}

// UpdateGeohash updates geohash prefix for spatial indexing
func (dl *DriverLocation) UpdateGeohash(geohashPrefix string) {
	dl.GeohashPrefix = geohashPrefix
	dl.UpdatedAt = time.Now()
}

// UpdateLocationHash updates location hash for duplicate detection
func (dl *DriverLocation) UpdateLocationHash(hash string) {
	dl.LastLocationHash = hash
}

// HasMovedSinceLastSync checks if driver has moved since last sync
func (dl *DriverLocation) HasMovedSinceLastSync(minDistanceMeters float64) bool {
	if dl.PreviousLocation == nil || dl.CurrentLocation == nil {
		return true
	}
	// Convert distance from meters to km
	distanceKm := dl.CurrentLocation.DistanceToKm(dl.PreviousLocation)
	return distanceKm*1000 >= minDistanceMeters
}
