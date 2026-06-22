package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// DriverStatus represents driver availability status
type DriverStatus string

const (
	DriverOffline DriverStatus = "offline"
	DriverOnline  DriverStatus = "online"
	DriverBusy    DriverStatus = "busy"
	DriverOnRide  DriverStatus = "on_ride"
)

// VehicleType represents vehicle category
type VehicleType string

const (
	VehicleEconomy VehicleType = "economy"
	VehicleComfort VehicleType = "comfort"
	VehiclePremium VehicleType = "premium"
	VehicleXL      VehicleType = "xl"
)

// DriverProfile represents driver-specific profile
type DriverProfile struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex"`
	LicenseNumber     string            `gorm:"not null;uniqueIndex"`
	LicenseExpiry     time.Time         `gorm:"not null"`
	Status            DriverStatus      `gorm:"type:driver_status;default:'offline'"`
	IsVerified        bool              `gorm:"default:false"`
	Rating            float64           `gorm:"default:0"`
	TotalRatings      int               `gorm:"default:0"`
	TotalRides        int               `gorm:"default:0"`
	AcceptanceRate    float64           `gorm:"default:100"`
	CancellationRate  float64           `gorm:"default:0"`
	TotalEarnings     decimal.Decimal   `gorm:"type:numeric(12,2);default:0"`
	LastOnlineAt      *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

// Vehicle represents a driver's vehicle
type Vehicle struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primaryKey"`
	DriverID           uuid.UUID       `gorm:"type:uuid;not null;index"`
	LicensePlate       string          `gorm:"not null;uniqueIndex"`
	VehicleType        VehicleType     `gorm:"not null"`
	Make               string
	Model              string
	Year               int
	Color              string
	VIN                string          `gorm:"uniqueIndex"`
	IsActive           bool            `gorm:"default:true"`
	RegistrationExpiry time.Time
	InsuranceExpiry    time.Time
	Odometer           int             // km
	Documents          datatypes.JSONMap `gorm:"type:jsonb"` // {registration: {...}, insurance: {...}}
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

// DriverDocument represents driver documents
type DriverDocument struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	DriverID        uuid.UUID `gorm:"type:uuid;not null;index"`
	DocumentType    string    // license, insurance, registration, vehicle_inspection
	DocumentURL     string
	ExpiryDate      *time.Time
	Status          string    // pending, approved, expired, rejected
	RejectionReason string
	UploadedAt      time.Time `gorm:"autoCreateTime"`
	ExpiresAt       time.Time
}

// DriverEarnings represents driver daily/weekly/monthly earnings
type DriverEarnings struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	DriverID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Period          string    // daily, weekly, monthly
	PeriodStart     time.Time
	PeriodEnd       time.Time
	RidesCompleted  int
	RidesCancelled  int
	TotalFare       decimal.Decimal `gorm:"type:numeric(12,2)"`
	Bonuses         decimal.Decimal `gorm:"type:numeric(12,2);default:0"`
	Deductions      decimal.Decimal `gorm:"type:numeric(12,2);default:0"`
	NetEarnings     decimal.Decimal `gorm:"type:numeric(12,2)"`
	AverageFare     decimal.Decimal `gorm:"type:numeric(10,2)"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

// DriverBankAccount represents driver bank details for payouts
type DriverBankAccount struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	DriverID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	BankName           string
	AccountNumber      string          `gorm:"not null"`
	AccountHolderName  string
	IFSC               string
	IsVerified         bool            `gorm:"default:false"`
	VerifiedAt         *time.Time
	Documents          datatypes.JSONMap `gorm:"type:jsonb"` // {bank_statement: {...}}
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

// DriverLocation represents real-time driver location
type DriverLocation struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	DriverID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Latitude           float64
	Longitude          float64
	Heading            int       // 0-359 degrees
	Speed              int       // km/h
	Accuracy           int       // meters
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies table name for DriverProfile
func (DriverProfile) TableName() string {
	return "drivers"
}

// TableName specifies table name for Vehicle
func (Vehicle) TableName() string {
	return "vehicles"
}

// TableName specifies table name for DriverDocument
func (DriverDocument) TableName() string {
	return "driver_documents"
}

// TableName specifies table name for DriverEarnings
func (DriverEarnings) TableName() string {
	return "driver_earnings"
}

// TableName specifies table name for DriverBankAccount
func (DriverBankAccount) TableName() string {
	return "driver_bank_accounts"
}

// TableName specifies table name for DriverLocation
func (DriverLocation) TableName() string {
	return "driver_locations"
}
