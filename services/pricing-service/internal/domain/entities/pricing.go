package entities

import "time"

// PricingRule defines pricing for a ride type in a city
type PricingRule struct {
	ID              string
	RideType        string    // ECONOMY, COMFORT, BUSINESS, POOL
	City            string
	BaseFare        float64   // Minimum charge
	DistanceRate    float64   // Per km
	TimeRate        float64   // Per minute
	MinimumFare     float64   // Absolute minimum
	SurgeFactorMax  float64   // Max surge multiplier
	TaxPercentage   float64   // Tax %
	PoolDiscount    float64   // Pooling discount %
	ActiveFrom      time.Time
	ActiveUntil     *time.Time
	Status          string    // ACTIVE, INACTIVE, ARCHIVED
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FareCalculation is the result of fare calculation
type FareCalculation struct {
	ID                   string
	RideID               string
	RideType             string
	DistanceMeters       int
	DurationSeconds      int
	PickupLat            float64
	PickupLng            float64
	DropoffLat           float64
	DropoffLng           float64
	BaseFare             float64
	DistanceFare         float64
	TimeFare             float64
	SubtotalBeforeSurge  float64
	SurgeMultiplier      float64
	SurgeAmount          float64
	Taxes                float64
	DiscountCodeID       *string
	DiscountAmount       float64
	FinalFare            float64
	IsPool               bool
	City                 string
	CalculatedAt         time.Time
	CreatedAt            time.Time
}

// SurgeHistory tracks surge multiplier over time
type SurgeHistory struct {
	ID               string
	Timestamp        time.Time
	City             string
	Latitude         float64
	Longitude        float64
	SurgeMultiplier  float64
	ActiveRides      int
	AvailableDrivers int
	Reason           string    // TIME_PEAK, EVENT, WEATHER, DEMAND, SUPPLY_DEMAND
	CreatedAt        time.Time
}

// DiscountCode represents a promotional code
type DiscountCode struct {
	ID                    string
	Code                  string
	DiscountType          string    // FIXED or PERCENTAGE
	DiscountValue         float64
	MaxDiscount           *float64  // For percentage: cap
	MinimumFareAmount     *float64  // Minimum ride fare
	MaxUses               *int
	UsesRemaining         int
	ValidFrom             time.Time
	ValidUntil            time.Time
	ApplicableToPooling   bool
	Status                string    // ACTIVE, EXPIRED, EXHAUSTED, DISABLED
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// PricingEvent represents a pricing operation audit log
type PricingEvent struct {
	ID              string
	EventType       string    // FARE_CALCULATED, SURGE_UPDATED, DISCOUNT_APPLIED
	RideID          *string
	DiscountCodeID  *string
	Details         map[string]interface{}
	CreatedAt       time.Time
}

// SurgeRequest contains surge calculation inputs
type SurgeRequest struct {
	City             string
	Latitude         float64
	Longitude        float64
	ActiveRides      int
	AvailableDrivers int
}

// SurgeMultiplierCalculation contains surge calculation output
type SurgeMultiplierCalculation struct {
	Multiplier       float64
	Reason           string
	ActiveRides      int
	AvailableDrivers int
	RatioDemandSupply float64
	TimeFactorMultiplier float64
	LocationFactorMultiplier float64
	CalculatedAt     time.Time
}
