package entities

import (
	"time"

	"github.com/google/uuid"
)

// PoolStatus represents pool state
type PoolStatus string

const (
	PoolStatusForming    PoolStatus = "FORMING"
	PoolStatusActive     PoolStatus = "ACTIVE"
	PoolStatusCompleted  PoolStatus = "COMPLETED"
	PoolStatusCancelled  PoolStatus = "CANCELLED"
)

// PoolGroup represents a ride pool (max 3 passengers)
type PoolGroup struct {
	ID                  string      `db:"id"`
	DriverID            string      `db:"driver_id"`
	Status              PoolStatus  `db:"status"`
	MaxSize             int         `db:"max_size"` // 2 or 3
	CurrentSize         int         `db:"current_size"`
	RideIDs             []string    // Will store as JSON in DB
	PassengerSequence   []string    // Pickup order
	DropoffSequence     []string    // Dropoff order
	TotalDistance       int         `db:"total_distance_meters"`
	TotalDuration       int         `db:"total_duration_seconds"`
	PooledFare          float64     `db:"pooled_fare"` // Total fare for pool
	IndividualFares     map[string]float64 // Per-rider breakdown
	CompatibilityScore  float64     `db:"compatibility_score"`
	EstimatedProfit     float64     `db:"estimated_profit"` // Platform profit
	CreatedAt           time.Time   `db:"created_at"`
	CompletedAt         *time.Time  `db:"completed_at"`
	UpdatedAt           time.Time   `db:"updated_at"`
}

// PoolRequest represents a ride eligible for pooling
type PoolRequest struct {
	RideID              string
	DriverID            string
	PickupLat           float64
	PickupLng           float64
	DropoffLat          float64
	DropoffLng          float64
	PickupAddress       string
	DropoffAddress      string
	EstimatedDistance   int
	EstimatedDuration   int
	EstimatedFare       float64
	FemaleOnly          bool
	MaxDetourMinutes    int // 10 default
	MaxWaitMinutes      int // 5 default
	MinRouteOverlap     float64 // 0.7 default (70%)
	CreatedAt           time.Time
}

// PoolCandidate represents a viable ride for pooling
type PoolCandidate struct {
	RideID              string
	PickupLat           float64
	PickupLng           float64
	DropoffLat          float64
	DropoffLng          float64
	PickupAddress       string
	DropoffAddress      string
	RouteOverlap        float64 // 0-1
	DetourMinutes       int
	WaitMinutes         int
	CompatibilityScore  float64
	SavingsPercentage   float64 // 20-30% typical
	Rank                int
}

// PoolRoute represents optimized route for pool
type PoolRoute struct {
	ID                  string
	PoolID              string
	PickupSequence      []string // [ride_id_1, ride_id_2, ride_id_3]
	DropoffSequence     []string
	TotalDistance       int      // meters
	TotalDuration       int      // seconds
	OptimalRoutePoints  []RoutePoint
	CreatedAt           time.Time
}

// RoutePoint represents a waypoint in optimized route
type RoutePoint struct {
	Latitude    float64
	Longitude   float64
	PointType   string // PICKUP, DROPOFF, WAYPOINT
	RideID      string // Which ride this point belongs to
	Sequence    int    // Order in route
}

// PoolCompatibility calculates score for two rides
type PoolCompatibility struct {
	RouteOverlapScore   float64 // 0-1 (weight: 0.4)
	ProfitabilityScore  float64 // 0-1 (weight: 0.3)
	ETASimilarityScore  float64 // 0-1 (weight: 0.2)
	PickupProximityScore float64 // 0-1 (weight: 0.1)
	FinalScore          float64 // 0-1 (weighted sum)
	IsViable            bool    // True if score > 0.5
}

// NewPoolGroup creates new pool
func NewPoolGroup(driverID string, maxSize int) *PoolGroup {
	return &PoolGroup{
		ID:               uuid.New().String(),
		DriverID:         driverID,
		Status:           PoolStatusForming,
		MaxSize:          maxSize,
		CurrentSize:      0,
		RideIDs:          []string{},
		PassengerSequence: []string{},
		DropoffSequence:  []string{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IndividualFares:  make(map[string]float64),
	}
}

// AddRide adds ride to pool
func (p *PoolGroup) AddRide(rideID string) bool {
	if p.CurrentSize >= p.MaxSize {
		return false
	}
	p.RideIDs = append(p.RideIDs, rideID)
	p.PassengerSequence = append(p.PassengerSequence, rideID)
	p.CurrentSize++
	p.UpdatedAt = time.Now()
	return true
}

// IsFull checks if pool is at max capacity
func (p *PoolGroup) IsFull() bool {
	return p.CurrentSize >= p.MaxSize
}

// CanBeActivated checks if pool ready to start
func (p *PoolGroup) CanBeActivated() bool {
	return p.CurrentSize >= 2 && p.Status == PoolStatusForming
}

// Activate transitions pool to ACTIVE
func (p *PoolGroup) Activate() {
	p.Status = PoolStatusActive
	p.UpdatedAt = time.Now()
}

// Complete transitions pool to COMPLETED
func (p *PoolGroup) Complete() {
	p.Status = PoolStatusCompleted
	now := time.Now()
	p.CompletedAt = &now
	p.UpdatedAt = now
}

// Cancel transitions pool to CANCELLED
func (p *PoolGroup) Cancel() {
	p.Status = PoolStatusCancelled
	now := time.Now()
	p.CompletedAt = &now
	p.UpdatedAt = now
}
